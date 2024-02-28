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
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer" // newLogsReceiver creates a new logs specific receiver.
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type logsGeneratorReceiver struct {
	telemetryGeneratorReceiver
	nextConsumer consumer.Logs
	logs         plog.Logs
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
	lr.generator = lr
	r.supportedTelemetry = component.DataTypeLogs

	lr.initializeLogs()

	return lr, nil
}

func (r *logsGeneratorReceiver) initializeLogs() {
	r.logs = plog.NewLogs()
	for _, g := range r.cfg.Generators {
		if g.Type != component.DataTypeLogs {
			continue
		}
		resourceLogs := r.logs.ResourceLogs().AppendEmpty()
		// Add resource attributes
		for k, v := range g.ResourceAttributes {
			resourceLogs.Resource().Attributes().PutStr(k, v)
		}
		scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
		// Generate logs
		logRecord := scopeLogs.LogRecords().AppendEmpty()
		for k, v := range g.Attributes {
			logRecord.Attributes().PutStr(k, v)
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
}

// generateTelemetry
func (r *logsGeneratorReceiver) generate() error {

	for i := 0; i < r.logs.ResourceLogs().Len(); i++ {
		resourceLogs := r.logs.ResourceLogs().At(i)
		for k := 0; k < resourceLogs.ScopeLogs().Len(); k++ {
			scopeLogs := resourceLogs.ScopeLogs().At(k)
			for j := 0; j < scopeLogs.LogRecords().Len(); j++ {
				log := scopeLogs.LogRecords().At(j)
				log.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
			}
		}
	}

	// Send logs to the next consumer
	return r.nextConsumer.ConsumeLogs(r.ctx, r.logs)
}
