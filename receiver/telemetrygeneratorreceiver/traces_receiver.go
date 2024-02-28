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

type tracesGeneratorReceiver struct {
	telemetryGeneratorReceiver
	nextConsumer consumer.Traces
}

// newTracesReceiver creates a new traces specific receiver.
func newTracesReceiver(ctx context.Context, id component.ID, logger *zap.Logger, cfg *Config, nextConsumer consumer.Traces) (*tracesGeneratorReceiver, error) {
	tr := &tracesGeneratorReceiver{
		nextConsumer: nextConsumer,
	}

	r, err := newTelemetryGeneratorReceiver(ctx, id, logger, cfg, tr)
	if err != nil {
		return nil, err
	}

	tr.telemetryGeneratorReceiver = r
	tr.generator = tr
	r.supportedTelemetry = component.DataTypeTraces

	return tr, nil
}

// TODO implement generate for metrics
func (r *tracesGeneratorReceiver) initializeMetrics() {
}

// TODO implement generate for traces
func (r *tracesGeneratorReceiver) generate() error {
	return nil
}
