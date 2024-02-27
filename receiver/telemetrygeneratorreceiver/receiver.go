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
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type telemetryGeneratorReceiver struct {
	logger             *zap.Logger
	id                 component.ID
	cfg                *Config
	supportedTelemetry component.DataType
	doneChan           chan struct{}
	ctx                context.Context
	cancelFunc         context.CancelCauseFunc
}

// newMetricsReceiver creates a new metrics specific receiver.
func newMetricsReceiver(id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Metrics) (*telemetryGeneratorReceiver, error) {
	r, err := newTelemetryGeneratorReceiver(id, logger, cfg)
	if err != nil {
		return nil, err
	}

	r.supportedTelemetry = component.DataTypeMetrics

	return r, nil
}

// newLogsReceiver creates a new logs specific receiver.
func newLogsReceiver(id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Logs) (*telemetryGeneratorReceiver, error) {
	r, err := newTelemetryGeneratorReceiver(id, logger, cfg)
	if err != nil {
		return nil, err
	}

	r.supportedTelemetry = component.DataTypeLogs

	return r, nil
}

// newTracesReceiver creates a new traces specific receiver.
func newTracesReceiver(id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Traces) (*telemetryGeneratorReceiver, error) {
	r, err := newTelemetryGeneratorReceiver(id, logger, cfg)
	if err != nil {
		return nil, err
	}

	r.supportedTelemetry = component.DataTypeTraces
	return r, nil
}

// newTelemetryGeneratorReceiver creates a new rehydration receiver
func newTelemetryGeneratorReceiver(id component.ID, logger *zap.Logger, cfg *Config) (*telemetryGeneratorReceiver, error) {

	ctx, cancel := context.WithCancelCause(context.Background())

	return &telemetryGeneratorReceiver{
		logger:     logger,
		id:         id,
		cfg:        cfg,
		doneChan:   make(chan struct{}),
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

// Start starts the telemetryGeneratorReceiver receiver
func (r *telemetryGeneratorReceiver) Start(ctx context.Context, host component.Host) error {

	go r.scrape()
	return nil
}

// Shutdown shuts down the rehydration receiver
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

// scrape
func (r *telemetryGeneratorReceiver) scrape() {
	defer close(r.doneChan)

	ticker := time.NewTicker(time.Second / time.Duration(r.cfg.PayloadsPerSecond))
	defer ticker.Stop()

	// Call once before the loop to ensure we do a collection before the first ticker
	r.generateTelemetry()
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			r.generateTelemetry()
		}
	}
}

// generateTelemetry
func (r *telemetryGeneratorReceiver) generateTelemetry() {

	return
}
