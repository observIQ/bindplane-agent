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

package telemetrygeneratorreceiver //import "github.com/observiq/bindplane-agent/receiver/telemetrygeneratorreceiver"

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	"go.uber.org/zap"
)

type generator interface {
	generate() error
}
type telemetryGeneratorReceiver struct {
	logger             *zap.Logger
	id                 component.ID
	cfg                *Config
	supportedTelemetry component.DataType
	doneChan           chan struct{}
	ctx                context.Context
	cancelFunc         context.CancelCauseFunc
	generator          generator
}
type logsGeneratorReceiver struct {
	telemetryGeneratorReceiver
	nextConsumer consumer.Logs
}
type metricsGeneratorReceiver struct {
	telemetryGeneratorReceiver
	nextConsumer consumer.Metrics
}

type tracesGeneratorReceiver struct {
	telemetryGeneratorReceiver
	nextConsumer consumer.Traces
}

// newMetricsReceiver creates a new metrics specific receiver.
func newMetricsReceiver(id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Metrics) (*metricsGeneratorReceiver, error) {
	mr := &metricsGeneratorReceiver{

		nextConsumer: nextConsumer,
	}
	r, err := newTelemetryGeneratorReceiver(id, logger, cfg, mr)
	if err != nil {
		return nil, err
	}

	mr.telemetryGeneratorReceiver = r
	mr.generator = mr
	r.supportedTelemetry = component.DataTypeMetrics

	return mr, nil
}

// newLogsReceiver creates a new logs specific receiver.
func newLogsReceiver(id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Logs) (*logsGeneratorReceiver, error) {
	lr := &logsGeneratorReceiver{
		nextConsumer: nextConsumer,
	}
	r, err := newTelemetryGeneratorReceiver(id, logger, cfg, lr)
	if err != nil {
		return nil, err
	}

	lr.telemetryGeneratorReceiver = r
	lr.generator = lr
	r.supportedTelemetry = component.DataTypeLogs

	return lr, nil
}

// newTracesReceiver creates a new traces specific receiver.
func newTracesReceiver(id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Traces) (*tracesGeneratorReceiver, error) {
	tr := &tracesGeneratorReceiver{
		nextConsumer: nextConsumer,
	}

	r, err := newTelemetryGeneratorReceiver(id, logger, cfg, tr)
	if err != nil {
		return nil, err
	}

	tr.telemetryGeneratorReceiver = r
	tr.generator = tr
	r.supportedTelemetry = component.DataTypeTraces

	return tr, nil
}

// newTelemetryGeneratorReceiver creates a new rehydration receiver
func newTelemetryGeneratorReceiver(id component.ID, logger *zap.Logger, cfg *Config, g generator) (telemetryGeneratorReceiver, error) {
	ctx, cancel := context.WithCancelCause(context.Background())

	return telemetryGeneratorReceiver{
		logger:     logger,
		id:         id,
		cfg:        cfg,
		doneChan:   make(chan struct{}),
		ctx:        ctx,
		cancelFunc: cancel,
		generator:  g,
	}, nil
}

// Shutdown shuts down the telemetry generator receiver
func (r *telemetryGeneratorReceiver) Shutdown(ctx context.Context) error {
	r.cancelFunc(errors.New("shutdown"))
	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-r.doneChan:
	}

	return err
}

// Start starts the logsGeneratorReceiver receiver
func (r *telemetryGeneratorReceiver) Start(ctx context.Context, _ component.Host) error {

	go func() {
		defer close(r.doneChan)

		ticker := time.NewTicker(time.Second / time.Duration(r.cfg.PayloadsPerSecond))
		defer ticker.Stop()

		// Call once before the loop to ensure we do a collection before the first ticker
		r.generator.generate()
		for {
			select {
			case <-ctx.Done():
				return
			case <-r.ctx.Done():
				return
			case <-ticker.C:
				r.generator.generate()
			}
		}
	}()
	return nil
}

// generateTelemetry
func (r *logsGeneratorReceiver) generate() error {

	// Loop through the generators and generate telemetry

	logs := plog.NewLogs()
	for _, g := range r.cfg.Generators {
		if g.Type != component.DataTypeLogs {
			continue
		}
		resourceLogs := logs.ResourceLogs().AppendEmpty()
		// Add resource attributes
		for k, v := range g.ResourceAttributes {
			resourceLogs.Resource().Attributes().PutStr(k, v)
		}
		scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
		// Generate logs
		logRecord := scopeLogs.LogRecords().AppendEmpty()
		for k, v := range g.Attributes {
			logRecord.Attributes().PutStr(k, v)
			logRecord.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
		}
		for k, v := range g.AdditionalConfig {
			switch k {
			case "body":
				// parses body string and sets that as log body, but uses string if parsing fails
				parsedBody := map[string]any{}
				if err := json.Unmarshal([]byte(v.(string)), &parsedBody); err != nil {
					r.logger.Warn("unable to unmarshal log body", zap.Error(err))
					logRecord.Body().SetStr(v.(string))
				} else {
					if err := logRecord.Body().SetEmptyMap().FromRaw(parsedBody); err != nil {
						r.logger.Warn("failed to set body to parsed value", zap.Error(err))
						logRecord.Body().SetStr(v.(string))
					}
				}
				logRecord.Body().SetStr(v.(string))
			case "severity":
				logRecord.SetSeverityNumber(plog.SeverityNumber(v.(int)))
			}
		}
	}
	// Send logs to the next consumer
	return r.nextConsumer.ConsumeLogs(r.ctx, logs)
}

// TODO implement generate for metrics
func (r *metricsGeneratorReceiver) generate() error {
	return nil
}

// TODO implement generate for traces
func (r *tracesGeneratorReceiver) generate() error {
	return nil
}
