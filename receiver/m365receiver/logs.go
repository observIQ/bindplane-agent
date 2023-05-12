// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package m365receiver // import "github.com/observiq/observiq-otel-collector/receiver/m365receiver"

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/extension/experimental/storage"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

const (
	logStorageKey = "last_recorded_event"
)

type lClient interface {
	GetJSON(endpoint string) ([]jsonLogs, error)
	GetToken() error
	StartSubscription(endpoint string) error
	shutdown() error
}

type m365LogsReceiver struct {
	settings      component.TelemetrySettings
	logger        *zap.Logger
	consumer      consumer.Logs
	cfg           *Config
	client        lClient
	storageClient storage.Client
	id            component.ID

	wg           *sync.WaitGroup
	mu           sync.Mutex
	consumerMu   sync.Mutex
	pollInterval time.Duration
	cancel       context.CancelFunc
	audits       []auditMetaData
	record       *logRecord
	root         string
	startRoot    string
}

type logRecord struct {
	NextStartTime *time.Time `mapstructure:"next_start_time"`
}

type auditMetaData struct {
	name    string
	route   string
	enabled bool
}

func newM365Logs(cfg *Config, settings receiver.CreateSettings, consumer consumer.Logs) *m365LogsReceiver {
	return &m365LogsReceiver{
		settings:      settings.TelemetrySettings,
		logger:        settings.Logger,
		consumer:      consumer,
		cfg:           cfg,
		storageClient: storage.NewNopClient(),
		id:            settings.ID,
		wg:            &sync.WaitGroup{},
		pollInterval:  cfg.Logs.PollInterval,
		audits: []auditMetaData{
			{"general", "Audit.General", cfg.Logs.GeneralLogs},
			{"exchange", "Audit.Exchange", cfg.Logs.ExchangeLogs},
			{"sharepoint", "Audit.SharePoint", cfg.Logs.SharepointLogs},
			{"azureAD", "Audit.AzureActiveDirectory", cfg.Logs.AzureADLogs},
			{"dlp", "DLP.All", cfg.Logs.DLPLogs},
		},
		root:      fmt.Sprintf("https://manage.office.com/api/v1.0/%s/activity/feed/subscriptions/content?contentType=", cfg.TenantID),
		startRoot: fmt.Sprintf("https://manage.office.com/api/v1.0/%s/activity/feed/subscriptions/start?contentType=", cfg.TenantID),
	}
}

// creates the client for http requests, and initializes metadata struct
func (l *m365LogsReceiver) Start(ctx context.Context, host component.Host) error {
	// create default client
	httpClient, err := l.cfg.ToClient(host, l.settings)
	if err != nil {
		l.logger.Error("error creating HTTP client", zap.Error(err))
		return err
	}

	// create m365 log client, create token and start audit subscriptions
	l.client = newM365Client(httpClient, l.cfg, "https://manage.office.com/.default")
	err = l.client.GetToken()
	if err != nil {
		l.logger.Error("error creating authorization token", zap.Error(err))
		return err
	}
	for _, a := range l.audits {
		err = l.client.StartSubscription(l.startRoot + a.route)
		if err != nil {
			l.logger.Error("error starting audit subscriptions", zap.Error(err))
			return err
		}
	}

	// set cancel function
	cancelCtx, cancel := context.WithCancel(ctx)
	l.cancel = cancel

	// init checkpoint stuff
	storageClient, err := adapter.GetStorageClient(ctx, host, l.cfg.StorageID, l.id)
	if err != nil {
		return fmt.Errorf("failed to get storage client: %w", err)
	}
	l.storageClient = storageClient
	l.loadCheckpoint(cancelCtx)

	return l.startPolling(cancelCtx)
}

func (l *m365LogsReceiver) Shutdown(ctx context.Context) error {
	l.logger.Debug("shutting down logs receiver")
	l.cancel()
	l.wg.Wait()
	return l.checkpoint(ctx)
}

// spins a go routine at each poll interval to go get logs
func (l *m365LogsReceiver) startPolling(ctx context.Context) error {
	t := time.NewTicker(l.pollInterval)
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			select {
			case <-t.C:
				if err := l.pollLogs(ctx); err != nil {
					l.logger.Error("error while polling for logs", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

// spins a go routine for each audit type/endpoint
func (l *m365LogsReceiver) pollLogs(ctx context.Context) error {
	st := pcommon.NewTimestampFromTime(time.Now().Add(-l.pollInterval)).AsTime()
	now := time.Now()

	auditWG := &sync.WaitGroup{}
	auditWG.Add(5)
	for i := 0; i < len(l.audits); i++ {
		endpoint := l.root + l.audits[i].route + fmt.Sprintf("&;startTime=%s&;endTime=%s", st, now)
		go l.poll(ctx, now, &l.audits[i], endpoint, auditWG)
	}
	auditWG.Wait()

	l.record.NextStartTime = &now
	return l.checkpoint(ctx)
}

// collects log data from endpoint, transforms logs, consumes logs
func (l *m365LogsReceiver) poll(ctx context.Context, now time.Time, audit *auditMetaData, endpoint string, wg *sync.WaitGroup) {
	defer wg.Done()
	if !audit.enabled {
		return
	}

	l.mu.Lock()
	logData, err := l.client.GetJSON(endpoint)
	if err != nil {
		if err.Error() == "authorization denied" { // troubleshoot stale token
			l.logger.Debug("possible stale token; attempting to regenerate")
			err = l.client.GetToken()
			if err != nil { // something went wrong generating token
				l.logger.Error("error creating authorization token", zap.Error(err))
				l.mu.Unlock()
				return
			}
			logData, err = l.client.GetJSON(endpoint)
			if err != nil { // not a stale token error, unsure what is wrong
				l.logger.Error("unable to retrieve logs", zap.Error(err))
				l.mu.Unlock()
				return
			}
		}
		l.logger.Error("error retrieving logs", zap.Error(err))
		l.mu.Unlock()
		return
	}
	l.mu.Unlock()

	logs := l.transformLogs(pcommon.NewTimestampFromTime(now), audit, logData)

	l.consumerMu.Lock()
	if logs.LogRecordCount() > 0 {
		if err = l.consumer.ConsumeLogs(ctx, logs); err != nil {
			l.logger.Error("error consuming events", zap.Error(err))
		}
	}
	l.consumerMu.Unlock()
}

// constructs logs from logData
func (l *m365LogsReceiver) transformLogs(now pcommon.Timestamp, audit *auditMetaData, logData []jsonLogs) plog.Logs {
	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	ra := resourceLogs.Resource().Attributes()
	ra.PutStr("m365.audit", audit.name)
	ra.PutStr("m365.organization_id", l.cfg.TenantID)

	for _, log := range logData {
		// log body
		logRecord := resourceLogs.ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()
		bodyBytes, err := json.Marshal(log)
		if err != nil {
			l.logger.Error("unable to marshal event into body string", zap.Error(err))
		}
		logRecord.Body().SetStr(string(bodyBytes))

		// timestamp
		const layout = "2006-01-02T15:04:05"
		ts, err := time.Parse(layout, log.CreationTime)
		if err != nil {
			l.logger.Warn("unable to interpret when an event was created, expecting a RFC3339 timestamp", zap.String("timestamp", log.CreationTime), zap.String("log", log.Id))
			logRecord.SetTimestamp(now)
		} else {
			logRecord.SetTimestamp(pcommon.NewTimestampFromTime(ts))
		}
		logRecord.SetObservedTimestamp(now)

		// attributes
		attrs := logRecord.Attributes()
		attrs.PutStr("id", log.Id)
		attrs.PutStr("operation", log.Operation)
		attrs.PutStr("user_id", log.UserId)
		attrs.PutStr("user_type", matchUserType(log.UserType))

		parseOptionalAttributes(&attrs, &log)
	}

	return logs
}

// sets the checkpoint
func (l *m365LogsReceiver) checkpoint(ctx context.Context) error {
	bytes, err := json.Marshal(l.record)
	if err != nil {
		return fmt.Errorf("unable to write checkpoint: %w", err)
	}
	return l.storageClient.Set(ctx, logStorageKey, bytes)
}

// loads the checkpoint
func (l *m365LogsReceiver) loadCheckpoint(ctx context.Context) {
	bytes, err := l.storageClient.Get(ctx, logStorageKey)
	if err != nil {
		l.logger.Info("unable to load checkpoint from storage client, continuing without a previous checkpoint", zap.Error(err))
		l.record = &logRecord{}
		return
	}

	if bytes == nil {
		l.record = &logRecord{}
	}

	var record logRecord
	if err = json.Unmarshal(bytes, &record); err != nil {
		l.logger.Error("unable to decode stored record for events, continuing without a checkpoint", zap.Error(err))
		l.record = &logRecord{}
		return
	}
	l.record = &record
}

// adds any attributes to log that may or may not be present
func parseOptionalAttributes(m *pcommon.Map, log *jsonLogs) {
	if log.Workload != "" {
		m.PutStr("workload", log.Workload)
	}
	if log.ResultStatus != "" {
		m.PutStr("result_status", log.ResultStatus)
	}
	if log.SharepointSite != "" {
		m.PutStr("sharepoint.site.id", log.SharepointSite)
	}
	if log.SharepointSourceFileName != "" {
		m.PutStr("sharepoint.source.file.name", log.SharepointSourceFileName)
	}
	if log.ExchangeMailboxGUID != "" {
		m.PutStr("exchange.mailbox.id", log.ExchangeMailboxGUID)
	}
	if log.AzureActor != nil {
		m.PutStr("azure.actor.id", log.AzureActor.ID)
		m.PutStr("azure.actor.type", matchAzureUserType(log.AzureActor.Type))
	}
	if log.DLPSharePointMetaData != nil {
		m.PutStr("dlp.sharepoint.user", log.DLPSharePointMetaData.From)
	}
	if log.DLPExchangeMetaData != nil {
		m.PutStr("dlp.exchange.message.id", log.DLPExchangeMetaData.MessageID)
	}
	if log.DLPPolicyDetails != nil {
		m.PutStr("dlp.policy_details.policy.id", log.DLPPolicyDetails.PolicyId)
		m.PutStr("dlp.policy_details.policy.name", log.DLPPolicyDetails.PolicyName)
	}
	if log.SecurityAlertId != "" {
		m.PutStr("security.alert.id", log.SecurityAlertId)
	}
	if log.SecurityAlertName != "" {
		m.PutStr("security.alert.name", log.SecurityAlertName)
	}
	if log.YammerActorId != "" {
		m.PutStr("yammer.user.id", log.YammerActorId)
	}
	if log.YammerFileId != nil {
		m.PutStr("yammer.file.id", strconv.Itoa(*log.YammerFileId))
	}
	if log.DefenderEmail != nil {
		m.PutStr("defender.email.attachment", log.DefenderEmail.FileName)
	}
	if log.DefenderURL != "" {
		m.PutStr("defender.url.url", log.DefenderURL)
	}
	if log.DefenderFile != nil {
		m.PutStr("defender.file.id", log.DefenderFile.DocumentId)
		m.PutStr("defender.file.verdict", matchFileVerdict(log.DefenderFile.FileVerdict))
	}
	if log.DefenderFileSource != nil {
		m.PutStr("defender.file.source_workload", matchSourceWorkload(*log.DefenderFileSource))
	}
	if log.InvestigationId != "" {
		m.PutStr("investigation.id", log.InvestigationId)
	}
	if log.InvestigationStatus != "" {
		m.PutStr("investigation.status", log.InvestigationStatus)
	}
	if log.PowerAppName != "" {
		m.PutStr("powerbi.app.name", log.PowerAppName)
	}
	if log.DynamicsEntityId != "" {
		m.PutStr("dynamics365.entity.id", log.DynamicsEntityId)
	}
	if log.DynamicsEntityName != "" {
		m.PutStr("dynamics365.entity.name", log.DynamicsEntityName)
	}
	if log.QuarantineSource != nil {
		m.PutStr("quarantine.request_source", matchQuarantineSource(*log.QuarantineSource))
	}
	if log.FormId != "" {
		m.PutStr("forms.form.id", log.FormId)
	}
	if log.MIPLabelId != "" {
		m.PutStr("mip.label.id", log.MIPLabelId)
	}
	if log.EncryptedMessageId != "" {
		m.PutStr("encrypted_portal.message.id", log.EncryptedMessageId)
	}
	if log.CommCompliance != nil {
		m.PutStr("communication_compliance.exchange.network_message_id", log.CommCompliance.NetworkMessageId)
	}
	if log.ConnectorJobId != "" {
		m.PutStr("compliance_connector.job.id", log.ConnectorJobId)
	}
	if log.ConnectorTaskId != "" {
		m.PutStr("compliance_connector.task.id", log.ConnectorTaskId)
	}
	if log.DataShareInvitation != nil {
		m.PutStr("system_sync.data_share.share.id", log.DataShareInvitation.ShareId)
	}
	if log.MSGraphConsentAppId != "" {
		m.PutStr("graph.consent.app.id", log.MSGraphConsentAppId)
	}
	if log.VivaGoalsUsername != "" {
		m.PutStr("viva_goals.username", log.VivaGoalsUsername)
	}
	if log.VivaGoalsOrgName != "" {
		m.PutStr("viva_goals.organization", log.VivaGoalsOrgName)
	}
	if log.MSToDoAppId != "" {
		m.PutStr("to_do.app.id", log.MSToDoAppId)
	}
	if log.MSToDoItemId != "" {
		m.PutStr("to_do.item.id", log.MSToDoItemId)
	}
	if log.MSWebProjectId != "" {
		m.PutStr("web_project.project.id", log.MSWebProjectId)
	}
	if log.MSWebRoadmapId != "" {
		m.PutStr("web_project.road_map.id", log.MSWebRoadmapId)
	}
	if log.MSWebRoadmapItemId != "" {
		m.PutStr("web_project.road_map.item.id", log.MSWebRoadmapItemId)
	}
}

func matchUserType(user int) string {
	switch user {
	case 0:
		return "Regular"
	case 1:
		return "Reserved"
	case 2:
		return "Admin"
	case 3:
		return "DcAdmin"
	case 4:
		return "System"
	case 5:
		return "Application"
	case 6:
		return "ServicePrincipal"
	case 7:
		return "CustomPolicy"
	case 8:
		return "SystemPolicy"
	default:
		return "impossible"
	}
}

func matchAzureUserType(user int) string {
	switch user {
	case 0:
		return "Claim"
	case 1:
		return "Name"
	case 2:
		return "Other"
	case 3:
		return "PUID"
	case 4:
		return "SPN"
	case 5:
		return "UPN"
	default:
		return "impossible"
	}
}

func matchFileVerdict(v int) string {
	switch v {
	case 0:
		return "Good"
	case 1:
		return "Bad"
	case -1:
		return "Error"
	case -2:
		return "Timeout"
	case -3:
		return "Pending"
	default:
		return "impossible"
	}
}

func matchSourceWorkload(w int) string {
	switch w {
	case 0:
		return "SharePoint Online"
	case 1:
		return "OneDrive for Business"
	case 2:
		return "Microsoft Teams"
	default:
		return "impossible"
	}
}

func matchQuarantineSource(s int) string {
	switch s {
	case 0:
		return "SCC"
	case 1:
		return "Cmdlet"
	case 2:
		return "URLlink"
	default:
		return "impossible"
	}
}
