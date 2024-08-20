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

package oktareceiver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type oktaLogsReceiver struct {
	cfg       Config
	client    httpClient
	consumer  consumer.Logs
	doneChan  chan bool
	logger    *zap.Logger
	nextUrl   string
	startTime time.Time
	wg        *sync.WaitGroup
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// newOktaLogsReceiver returns a newly configured oktaLogsReceiver
func newOktaLogsReceiver(cfg *Config, logger *zap.Logger, consumer consumer.Logs) (*oktaLogsReceiver, error) {
	var startTime time.Time
	if cfg.StartTime != "" {
		cfgStartTime, err := time.Parse(OktaTimeFormat, cfg.StartTime)
		if err != nil {
			return nil, err
		}
		startTime = cfgStartTime
	} else {
		startTime = time.Now().UTC().Add(-cfg.PollInterval)
	}

	return &oktaLogsReceiver{
		cfg:       *cfg,
		client:    http.DefaultClient,
		consumer:  consumer,
		doneChan:  make(chan bool),
		logger:    logger,
		startTime: startTime,
		wg:        &sync.WaitGroup{},
	}, nil
}

func (r *oktaLogsReceiver) Start(ctx context.Context, host component.Host) error {
	r.wg.Add(1)
	go r.startPolling(ctx)
	return nil
}

func (r *oktaLogsReceiver) startPolling(ctx context.Context) {
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
		case <-r.doneChan:
			return
		case <-t.C:
			err := r.poll(ctx)
			if err != nil {
				r.logger.Error("there was an error during the poll", zap.Error(err))
			}
		}
	}
}

func (r *oktaLogsReceiver) poll(ctx context.Context) error {
	select {
	case _, ok := <-r.doneChan:
		if !ok {
			return nil
		}
	default:
		logEvents := r.getLogs()
		fmt.Println("\033[32m"+time.Now().Format("2006-01-02 15:04:05")+":", "Received", len(logEvents), "logs"+"\033[0m")
		observedTime := pcommon.NewTimestampFromTime(time.Now())
		logs := r.processLogEvents(observedTime, logEvents)
		if logs.LogRecordCount() > 0 {
			if err := r.consumer.ConsumeLogs(ctx, logs); err != nil {
				r.logger.Error("unable to consume logs", zap.Error(err))
				break
			}
		}
	}
	return nil
}

func (r *oktaLogsReceiver) getLogs() []*okta.LogEvent {
	var req *http.Request
	var err error
	var logs []*okta.LogEvent
	if r.nextUrl == "" {
		fmt.Println("\033[32m" + "first poll" + "\033[0m")
		// for the first polling request, use startTime
		reqUrl := "https://" + r.cfg.Domain + "/api/v1/logs"
		req, err = http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			r.logger.Warn("error creating okta api request", zap.Error(err))
			return logs
		}

		query := req.URL.Query()
		query.Add("since", r.startTime.Format(OktaTimeFormat))
		req.URL.RawQuery = query.Encode()
		fmt.Println("\033[32m" + "since: " + r.startTime.Format(OktaTimeFormat) + "\033[0m")
	} else {
		req, err = http.NewRequest("GET", r.nextUrl, nil)
		if err != nil {
			r.logger.Warn("error creating okta api request", zap.Error(err))
			return logs
		}
	}

	req.Header.Add("Authorization", "SSWS "+r.cfg.ApiToken)
	res, err := r.client.Do(req)
	if err != nil {
		r.logger.Warn("error performing okta api request", zap.Error(err))
		return logs
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		r.logger.Error("okta logs endpoint returned non-200 statuscode: " + res.Status)
		return logs
	}

	r.setNextLink(res)
	fmt.Println("\033[32m"+"Next Url:", r.nextUrl+"\033[0m")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		r.logger.Warn("error reading response body", zap.Error(err))
		return logs
	}

	err = json.Unmarshal(body, &logs)
	if err != nil {
		r.logger.Warn("unable to unmarshal log event", zap.Error(err))
	}

	return logs
}

func (r *oktaLogsReceiver) processLogEvents(observedTime pcommon.Timestamp, logEvents []*okta.LogEvent) plog.Logs {
	logs := plog.NewLogs()

	for _, logEvent := range logEvents {
		resourceLogs := logs.ResourceLogs().AppendEmpty()
		resourceLogs.ScopeLogs().AppendEmpty()
		logRecord := resourceLogs.ScopeLogs().At(0).LogRecords().AppendEmpty()

		// resource attributes
		resourceAttributes := resourceLogs.Resource().Attributes()
		resourceAttributes.PutStr("okta.domain", r.cfg.Domain)

		// timestamps
		logRecord.SetObservedTimestamp(observedTime)
		timestamp := time.UnixMilli(logEvent.Published.UnixMilli())
		logRecord.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))

		// body
		logEventBytes, err := json.Marshal(logEvent)
		if err != nil {
			r.logger.Warn("unable to marshal logEvent", zap.Error(err))
		} else {
			logRecord.Body().SetStr(string(logEventBytes))
		}

		// attributes
		logRecord.Attributes().PutStr("uuid", logEvent.Uuid)
		logRecord.Attributes().PutStr("eventType", logEvent.EventType)
		logRecord.Attributes().PutStr("displayMessage", logEvent.DisplayMessage)
		logRecord.Attributes().PutStr("outcome.result", logEvent.Outcome.Result)
		logRecord.Attributes().PutStr("actor.id", logEvent.Actor.Id)
		logRecord.Attributes().PutStr("actor.alternateId", logEvent.Actor.AlternateId)
		logRecord.Attributes().PutStr("actor.displayName", logEvent.Actor.DisplayName)
	}

	return logs
}

func (r *oktaLogsReceiver) Shutdown(_ context.Context) error {
	r.logger.Debug("shutting down logs receiver")
	close(r.doneChan)
	r.wg.Wait()
	return nil
}

func (r *oktaLogsReceiver) setNextLink(res *http.Response) {
	for _, link := range res.Header["Link"] {
		// Split the link into URL and parameters
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) < 2 {
			continue
		}

		// Check if the "rel" parameter is "next"
		if strings.TrimSpace(parts[1]) == `rel="next"` {
			// Extract and return the URL
			r.nextUrl = strings.Trim(parts[0], "<>")
			return
		}
	}
	r.logger.Warn("unable to get next link")
}
