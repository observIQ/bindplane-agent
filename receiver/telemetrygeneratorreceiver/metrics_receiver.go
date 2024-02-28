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

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

type metricsGeneratorReceiver struct {
	telemetryGeneratorReceiver
	nextConsumer consumer.Metrics
}

// newMetricsReceiver creates a new metrics specific receiver.
func newMetricsReceiver(ctx context.Context, logger *zap.Logger, cfg *Config, nextConsumer consumer.Metrics) (*metricsGeneratorReceiver, error) {
	mr := &metricsGeneratorReceiver{

		nextConsumer: nextConsumer,
	}
	r, err := newTelemetryGeneratorReceiver(ctx, logger, cfg, mr)
	if err != nil {
		return nil, err
	}

	mr.telemetryGeneratorReceiver = r
	mr.generator = mr
	r.supportedTelemetry = component.DataTypeMetrics

	mr.initializeMetrics()

	return mr, nil
}

// TODO implement
func (r *metricsGeneratorReceiver) initializeMetrics() {
}

// TODO implement generate for metrics
func (r *metricsGeneratorReceiver) generate() error {
	return nil
}
