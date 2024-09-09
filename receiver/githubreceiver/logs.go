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

package githubreceiver // import "github.com/observiq/bindplane-agent/receiver/github"

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

var gitHubMaxLimit = 100




type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

// githubLogsReceiver is a receiver for GitHub logs.
type gitHubLogsReceiver struct {
	client   httpClient
	logger   *zap.Logger
	consumer consumer.Logs
	cfg      *Config
	wg       *sync.WaitGroup
	cancel   context.CancelFunc
	nextURL  string
}

// newGitHubLogsReceiver creates a new GitHub logs receiver.
func newGitHubLogsReceiver(cfg *Config, logger *zap.Logger, consumer consumer.Logs) (*gitHubLogsReceiver, error) {
	return &gitHubLogsReceiver{
		cfg:      cfg,
		logger:   logger,
		consumer: consumer,
		wg:       &sync.WaitGroup{},
		client:   http.DefaultClient,
	}, nil
}

// Start begins the receiver's operation.
func (r *gitHubLogsReceiver) Start(_ context.Context, _ component.Host) error {
	r.logger.Info("Starting GitHub logs receiver")
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.wg.Add(1)
	if r.cfg.WebhookConfig != nil {
		// Set up webhook handling
		go r.setupWebhook()
	}
	if r.cfg.PollInterval > 0 {
		// Start polling
		go r.startPolling(ctx)
	}
	return nil

}

func (r *gitHubLogsReceiver) setupWebhook() error {
	// waiting on open source code for webhooks
	return nil
}

func (r *gitHubLogsReceiver) startPolling(ctx context.Context) {
	defer r.wg.Done()
	t := time.NewTicker(r.cfg.PollInterval)
	err := r.poll(ctx)
	if err != nil {
		r.logger.Error("there was an error during the first poll", zap.Error(err))
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			err := r.poll(ctx)
			if err != nil {
				r.logger.Error("there was an error during the poll", zap.Error(err))
			}
		}
	}
}

func (r *gitHubLogsReceiver) poll(ctx context.Context) error {
	logEvents := r.getLogs(ctx)
	observedTime := pcommon.NewTimestampFromTime(time.Now())
	logs := r.processLogEvents(observedTime, logEvents)
	if logs.LogRecordCount() > 0 {
		if err := r.consumer.ConsumeLogs(ctx, logs); err != nil {
			return err
		}
	}
	return nil
}

func (r *gitHubLogsReceiver) getLogs(ctx context.Context) ([]gitHubEnterpriseLog) {
	pollTime := time.Now().UTC()
	token := r.cfg.AccessToken
	var endpoint string
	var curLogs []gitHubEnterpriseLog
	var allLogs []gitHubEnterpriseLog
	var url string
	page := 1
	// endpoint changes based on log type
	if r.cfg.LogType == "user" {
		endpoint = fmt.Sprintf("users/%s/events/public", r.cfg.Name)
	} else if r.cfg.LogType == "organization" {
		endpoint = fmt.Sprintf("orgs/%s/audit-log", r.cfg.Name)
	} else {

		endpoint = fmt.Sprintf("enterprises/%s/audit-log", r.cfg.Name)
	}	
	// Set the initial URL


	for {
		// Use nextURL if it's set
		if r.nextURL != "" {
			url = r.nextURL
		} else {
			url = fmt.Sprintf("https://api.github.com/%s?per_page=%d&page=%d", endpoint, gitHubMaxLimit, page) 
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			r.logger.Error("error creating request: %w", zap.Error(err))
			return allLogs
		}

		// Add headers
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Accept", "application/vnd.github.v3+json") // optional?
		req.Header.Add("X-GitHub-Api-Version", "2022-11-28") // optional?
		
		resp, err := r.client.Do(req)
		if err != nil {
			r.logger.Error("error making request", zap.Error(err))
			return allLogs
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			r.logger.Error("unexpected status code", zap.Int("statusCode", resp.StatusCode))
			return allLogs
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			r.logger.Error("error reading response body", zap.Error(err))
			return allLogs
		}
		
		// Unmarshal JSON into curLogs
		err = json.Unmarshal(body, &curLogs)
		if err != nil {
			r.logger.Error("error unmarshalling JSON", zap.Error(err))
			return allLogs
		}
		
		// Append the current logs to the allLogs slice
		allLogs = append(allLogs, curLogs...)
		page = r.setNextLink(resp, page)
		// Check for the 'Link' header and set the next URL if it exists
		var curLogTime time.Time
		if len(curLogs) != 0 {
			curLogTime = millisecondsToTime(curLogs[len(curLogs)-1].Timestamp)
		}
		if r.nextURL == "" || len(curLogs) < gitHubMaxLimit || curLogTime.After(pollTime) {
			break
		}
	}
	return allLogs
}



func (r *gitHubLogsReceiver) processLogEvents(observedTime pcommon.Timestamp, logEvents []gitHubEnterpriseLog) plog.Logs {
	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	resourceLogs.ScopeLogs().AppendEmpty()
	resourceLogs.Resource().Attributes().PutStr("access_token", r.cfg.AccessToken)
	for _, logEvent := range logEvents {
		logRecord := resourceLogs.ScopeLogs().At(0).LogRecords().AppendEmpty()

		// timestamps
		logRecord.SetObservedTimestamp(observedTime)
		timestamp := time.UnixMilli(logEvent.Timestamp)
		logRecord.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))

		// body
		logEventBytes, err := json.Marshal(logEvent)
		if err != nil {
			r.logger.Error("unable to marshal logEvent", zap.Error(err))
		} else {
			logRecord.Body().SetStr(string(logEventBytes))
		}
		// add attributes
		r.addAttributes(logRecord, logEvent)
	}
	return logs
}

// helper function to add attributes to logRecord
func (r *gitHubLogsReceiver) addAttributes(logRecord plog.LogRecord, logEvent gitHubEnterpriseLog) {
	// add attributes
	attrs := logRecord.Attributes()
	attrs.PutInt("@timestamp", logEvent.Timestamp)
	attrs.PutStr("_document_id", logEvent.DocumentID)
	attrs.PutStr("action", logEvent.Action)
	attrs.PutStr("actor", logEvent.Actor)
	attrs.PutInt("actor_id", logEvent.ActorID)
	attrs.PutBool("actor_is_bot", logEvent.ActorIsBot)
	attrs.PutStr("actor_location", logEvent.ActorLocation.CountryCode)
	attrs.PutStr("business", logEvent.Business)
	attrs.PutInt("business_id", logEvent.BusinessID)
	attrs.PutInt("created_at", logEvent.CreatedAt)
	attrs.PutStr("operation_type", logEvent.OperationType)
	attrs.PutStr("user_agent", logEvent.UserAgent)

	// map of optional attributes
	optionalAttrs := map[string]interface{}{
		"name":                    logEvent.Name,
		"org":                     logEvent.Org,
		"org_id":                  logEvent.OrgID,
		"organization_upgrade":    logEvent.OrganizationUpgrade,
		"permission":              logEvent.Permission,
		"user":                    logEvent.User,
		"user_id":                 logEvent.UserID,
		"owner_type":              logEvent.OwnerType,
		"audit_log_stream_sink":   logEvent.AuditLogStreamSink,
		"audit_log_stream_result": logEvent.AuditLogStreamResult,
	}

	// add optional attributes if they are present
	for key, value := range optionalAttrs {
		switch v := value.(type) {
		case string:
			if v != "" {
				attrs.PutStr(key, v)
			}
		case int64:
			if v != 0 {
				attrs.PutInt(key, v)
			}
		case bool:
			if v {
				attrs.PutBool(key, v)
			}
		}
	}

}
func (r *gitHubLogsReceiver) setNextLink(res *http.Response, page int) int {
	for _, link := range res.Header["Link"] {
		// Split the link into URL and parameters
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) < 2 {
			continue
		}
		
		// Check if the "rel" parameter is "next"
		if strings.TrimSpace(parts[1]) == `rel="next"` {
			// increment page
			page ++
			// Extract and return the URL
			r.nextURL = strings.Trim(parts[0], "<>")
			
			return page
		}
	}
	r.logger.Error("unable to get next link")
	return page
}

func millisecondsToTime(ms int64) time.Time {
    seconds := ms / 1000
    nanoseconds := (ms % 1000) * 1000000
    return time.Unix(seconds, nanoseconds)
}

func (r *gitHubLogsReceiver) Shutdown(ctx context.Context) error {
	r.logger.Debug("shutting down logs receiver")
	if r.cancel != nil {
		r.cancel()
	}
	r.client.CloseIdleConnections()
	r.wg.Wait()
	return nil
}
