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
	"strconv"
	"sync"
	"time"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type oktaLogsReceiver struct {
	cfg           Config
	consumer      consumer.Logs
	wg            *sync.WaitGroup
	logger        *zap.Logger
	doneChan      chan bool
	pollInterval  time.Duration
	nextStartTime time.Time
}

// newOktaLogsReceiver returns a newly configured oktaLogsReceiver
func newOktaLogsReceiver(cfg *Config, logger *zap.Logger, consumer consumer.Logs) (*oktaLogsReceiver, error) {
	return &oktaLogsReceiver{
		cfg:           *cfg,
		consumer:      consumer,
		wg:            &sync.WaitGroup{},
		logger:        logger,
		doneChan:      make(chan bool),
		pollInterval:  cfg.PollInterval,
		nextStartTime: time.Now().Add(-cfg.PollInterval),
	}, nil
}

func (r *oktaLogsReceiver) Start(ctx context.Context, host component.Host) error {
	r.logger.Debug("starting to poll for Okta logs")
	r.wg.Add(1)
	go r.startPolling(ctx)
	return nil
}

func (r *oktaLogsReceiver) startPolling(ctx context.Context) {
	defer r.wg.Done()

	fmt.Println(r.cfg.PollInterval)
	t := time.NewTicker(r.cfg.PollInterval)
	r.poll(ctx)
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
	startTime := r.nextStartTime
	endTime := time.Now()
	err := r.pollForLogs(ctx, startTime, endTime)
	r.nextStartTime = endTime
	return err
}

func (r *oktaLogsReceiver) pollForLogs(ctx context.Context, startTime, endTime time.Time) error {
	select {
	// if done, we want to stop processing paginated stream of events
	case _, ok := <-r.doneChan:
		if !ok {
			return nil
		}
	default:
		logEvents := r.requestLogs()
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

func (r *oktaLogsReceiver) requestLogs() []*okta.LogEvent {
	ctx, client, err := okta.NewClient(
		context.TODO(),
		okta.WithOrgUrl("https://"+r.cfg.Domain),
		okta.WithToken(r.cfg.ApiToken),
	)
	if err != nil {
		panic(err)
	}

	qp := &query.Params{}
	qp.Limit = 1
	logs, resp, err := client.LogEvent.GetLogs(ctx, qp)
	if resp.StatusCode != 200 {
		r.logger.Error("okta sdk GetLogs returned non-200 statuscode: " + strconv.Itoa(resp.StatusCode))
	}
	if err != nil {
		panic(err)
	}
	return logs
}

func (r *oktaLogsReceiver) processLogEvents(observedTime pcommon.Timestamp, logEvents []*okta.LogEvent) plog.Logs {
	logs := plog.NewLogs()

	for _, logEvent := range logEvents {
		resourceLogs := logs.ResourceLogs().AppendEmpty()
		resourceLogs.ScopeLogs().AppendEmpty()
		resourceAttributes := resourceLogs.Resource().Attributes()
		resourceAttributes.PutStr("okta.domain", r.cfg.Domain)

		logRecord := resourceLogs.ScopeLogs().At(0).LogRecords().AppendEmpty()
		logRecord.SetObservedTimestamp(observedTime)
		timestamp := time.UnixMilli(logEvent.Published.UnixMilli())
		logRecord.SetTimestamp(pcommon.NewTimestampFromTime(timestamp))
		logEventBytes, err := json.Marshal(logEvent)
		if err != nil {
			panic(err)
		}
		logRecord.Body().SetStr(string(logEventBytes))
		logRecord.Attributes().PutStr("uuid", logEvent.Uuid)
		logRecord.Attributes().PutStr("event.type", logEvent.EventType)
	}

	return logs
}

func (r *oktaLogsReceiver) Shutdown(_ context.Context) error {
	r.logger.Debug("shutting down logs receiver")
	close(r.doneChan)
	r.wg.Wait()
	return nil
}
