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
	"encoding/json"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type logsGenerator struct {
	cfg    *GeneratorConfig
	logs   plog.Logs
	logger *zap.Logger
}

func newLogsGenerator(cfg *GeneratorConfig, logger *zap.Logger) generator {
	return &logsGenerator{
		cfg:    cfg,
		logger: logger,
		logs:   plog.NewLogs(),
	}
}

func (g *logsGenerator) initialize() {
	resourceLogs := g.logs.ResourceLogs().AppendEmpty()
	// Add resource attributes
	for k, v := range g.cfg.ResourceAttributes {
		resourceLogs.Resource().Attributes().PutStr(k, v)
	}
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
	// Generate logs
	logRecord := scopeLogs.LogRecords().AppendEmpty()
	for k, v := range g.cfg.Attributes {
		logRecord.Attributes().PutStr(k, v)
	}
	for k, v := range g.cfg.AdditionalConfig {
		switch k {
		case "body":
			// parses body string and sets that as log body, but uses string if parsing fails
			parsedBody := map[string]any{}
			if err := json.Unmarshal([]byte(v.(string)), &parsedBody); err != nil {
				g.logger.Warn("unable to unmarshal log body", zap.Error(err))
				logRecord.Body().SetStr(v.(string))
			} else {
				if err := logRecord.Body().SetEmptyMap().FromRaw(parsedBody); err != nil {
					g.logger.Warn("failed to set body to parsed value", zap.Error(err))
					logRecord.Body().SetStr(v.(string))
				}
			}
			logRecord.Body().SetStr(v.(string))
		case "severity":
			logRecord.SetSeverityNumber(plog.SeverityNumber(v.(int)))
		}
	}

}

func (g *logsGenerator) SupportsType(t component.Type) bool {
	return t == component.DataTypeLogs
}

func (g *logsGenerator) GenerateMetrics() pmetric.Metrics {
	return pmetric.NewMetrics()
}

func (g *logsGenerator) GenerateLogs() plog.Logs {
	for i := 0; i < g.logs.ResourceLogs().Len(); i++ {
		resourceLogs := g.logs.ResourceLogs().At(i)
		for k := 0; k < resourceLogs.ScopeLogs().Len(); k++ {
			scopeLogs := resourceLogs.ScopeLogs().At(k)
			for j := 0; j < scopeLogs.LogRecords().Len(); j++ {
				log := scopeLogs.LogRecords().At(j)
				log.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
			}
		}
	}
	return g.logs
}

func (g *logsGenerator) GenerateTraces() ptrace.Traces {
	return ptrace.NewTraces()
}
