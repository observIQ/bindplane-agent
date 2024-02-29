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
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type hostMetricsGenerator struct {
	cfg    GeneratorConfig
	logger *zap.Logger
}

func newHostMetricsGenerator(cfg GeneratorConfig, logger *zap.Logger) generator {
	return &hostMetricsGenerator{
		cfg:    cfg,
		logger: logger,
	}
}

func (g *hostMetricsGenerator) initialize() {

}

func (g *hostMetricsGenerator) SupportsType(t component.Type) bool {
	return t == component.DataTypeMetrics
}

func (g *hostMetricsGenerator) GenerateMetrics() pmetric.Metrics {
	return pmetric.NewMetrics()
}

func (g *hostMetricsGenerator) GenerateLogs() plog.Logs {
	return plog.NewLogs()
}

func (g *hostMetricsGenerator) GenerateTraces() ptrace.Traces {
	return ptrace.NewTraces()
}
