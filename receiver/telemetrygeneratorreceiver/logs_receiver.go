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
	"go.opentelemetry.io/collector/consumer" // newLogsReceiver creates a new logs specific receiver.
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type logsGeneratorReceiver struct {
	telemetryGeneratorReceiver
	nextConsumer consumer.Logs
	generators   []logGenerator
}

func newLogsReceiver(ctx context.Context, logger *zap.Logger, cfg *Config, nextConsumer consumer.Logs) (*logsGeneratorReceiver, error) {
	lr := &logsGeneratorReceiver{
		nextConsumer: nextConsumer,
	}
	r, err := newTelemetryGeneratorReceiver(ctx, logger, cfg, lr)
	if err != nil {
		return nil, err
	}

	lr.telemetryGeneratorReceiver = r
	r.supportedTelemetry = component.DataTypeLogs

	lr.generators = newLogsGenerators(cfg, logger)

	return lr, nil
}

// generateTelemetry
func (r *logsGeneratorReceiver) produce() error {
	logs := plog.NewLogs()
	for _, g := range r.generators {
		l := g.generateLogs()
		for i := 0; i < l.ResourceLogs().Len(); i++ {
			src := l.ResourceLogs().At(i)
			src.CopyTo(logs.ResourceLogs().AppendEmpty())
		}
	}
	return r.nextConsumer.ConsumeLogs(r.ctx, logs)
}
