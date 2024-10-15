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

package sentinelonereceiver

import (
	"context"
	"net/http"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type sentinelOneLogsReceiver struct {
	cfg      Config
	client   httpClient
	consumer consumer.Logs
	logger   *zap.Logger
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
}

type sentinelOneLogEvent struct{}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

// newSentinelOneLogsReceiver returns a newly configured sentinelOneLogsReceiver
func newSentinelOneLogsReceiver(cfg *Config, logger *zap.Logger, consumer consumer.Logs) (*sentinelOneLogsReceiver, error) {
	return &sentinelOneLogsReceiver{
		cfg:      *cfg,
		client:   http.DefaultClient,
		consumer: consumer,
		logger:   logger,
		wg:       &sync.WaitGroup{},
	}, nil
}

func (r *sentinelOneLogsReceiver) Start(_ context.Context, _ component.Host) error {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.wg.Add(1)
	go r.startPolling(ctx)
	return nil
}

func (r *sentinelOneLogsReceiver) startPolling(ctx context.Context) {
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

func (r *sentinelOneLogsReceiver) poll(ctx context.Context) error {
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

func (r *sentinelOneLogsReceiver) getLogs(ctx context.Context) []sentinelOneLogEvent {
	return []sentinelOneLogEvent{}
}

func (r *sentinelOneLogsReceiver) processLogEvents(observedTime pcommon.Timestamp, logEvents []sentinelOneLogEvent) plog.Logs {
	return plog.NewLogs()
}

func (r *sentinelOneLogsReceiver) Shutdown(_ context.Context) error {
	r.logger.Debug("shutting down logs receiver")
	if r.cancel != nil {
		r.cancel()
	}
	r.client.CloseIdleConnections()
	r.wg.Wait()
	return nil
}
