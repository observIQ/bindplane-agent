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
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/extension/experimental/storage"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

type lClient interface {
	GetJSON(endpoint string) ([]jsonLogs, error)
	GetToken() error
	shutdown() error
}

type m365LogsReceiver struct {
	settings      component.TelemetrySettings
	logger        *zap.Logger
	consumer      consumer.Logs
	cfg           *Config
	client        lClient
	storageClient storage.Client

	wg           *sync.WaitGroup
	mu           sync.Mutex
	consumerMu   sync.Mutex
	pollInterval time.Duration
	cancel       context.CancelFunc
	audits       []auditMetaData
	root         string
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
		wg:            &sync.WaitGroup{},
		pollInterval:  cfg.Logs.PollInterval,
		audits: []auditMetaData{
			{"general", "Audit.General", cfg.Logs.GeneralLogs},
			{"exchange", "Audit.Exchange", cfg.Logs.ExchangeLogs},
			{"sharepoint", "Audit.SharePoint", cfg.Logs.SharepointLogs},
			{"azureAD", "Audit.AzureActiveDirectory", cfg.Logs.AzureADLogs},
			{"dlp", "DLP.All", cfg.Logs.DLPLogs},
		},
		root: fmt.Sprintf("https://manage.office.com/api/v1.0/%s/activity/feed/subscriptions/content?contentType=", cfg.TenantID),
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

	// create m365 log client
	l.client = newM365Client(httpClient, l.cfg, "https://manage.office.com/.default")
	err = l.client.GetToken()
	if err != nil {
		l.logger.Error("error creating authorization token", zap.Error(err))
		return err
	}

	// set cancel function
	//TODO: add checkpoint stuff here when that point comes
	cancelCtx, cancel := context.WithCancel(ctx)
	l.cancel = cancel

	return l.startPolling(cancelCtx)
}

func (l *m365LogsReceiver) Shutdown(_ context.Context) error {
	l.logger.Debug("shutting down logs receiver")
	l.cancel()
	l.wg.Wait()
	return nil
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

	for _, a := range l.audits {
		endpoint := l.root + a.route + fmt.Sprintf("&;startTime=%s&;endTime=%s", st, now)
		l.wg.Add(1)
		go l.poll(ctx, now, &a, endpoint)
	}

	return nil
}

// collects log data from endpoint, transforms logs, consumes logs
func (l *m365LogsReceiver) poll(ctx context.Context, now time.Time, audit *auditMetaData, endpoint string) {
	defer l.wg.Done()
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

// adds any attributes to log that may or may not be present
func parseOptionalAttributes(m *pcommon.Map, log *jsonLogs) {
	if log.Workload != "" {
		m.PutStr("workload", log.Workload)
	}
	if log.ResultStatus != "" {
		m.PutStr("result_status", log.ResultStatus)
	}
}

func matchUserType(userType int) string {
	switch userType {
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
