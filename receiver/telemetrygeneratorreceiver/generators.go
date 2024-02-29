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

// GeneratorType is the type of generator to use, either "logs", "host_metrics", or "windows_events"
type GeneratorType = string

const (
	// GeneratorTypeLogs is the generator type for logs
	GeneratorTypeLogs GeneratorType = "logs"

	// GeneratorTypeHostMetrics is the generator type for host metrics
	GeneratorTypeHostMetrics GeneratorType = "host_metrics"

	// GeneratorTypeWindowsEvents is the generator type for windows events
	GeneratorTypeWindowsEvents GeneratorType = "windows_events"
)

type generator interface {
	// SupportsType returns true if the generator supports the given component type, either metrics, logs, or traces.
	SupportsType(component.Type) bool
	GenerateMetrics() pmetric.Metrics
	GenerateLogs() plog.Logs
	GenerateTraces() ptrace.Traces

	// initialize must called when the generator is created. It is used to prevent the need to
	// recreate the telemetry every generation cycle
	initialize()
}

func newGenerators(cfg *Config, logger *zap.Logger, supportedType component.DataType) []generator {
	var generators []generator
	for _, gen := range cfg.Generators {
		var newGenerator generator
		switch gen.Type {
		case GeneratorTypeLogs:
			newGenerator = newLogsGenerator(gen, logger)
		case GeneratorTypeHostMetrics:
			newGenerator = newHostMetricsGenerator(gen, logger)
		case GeneratorTypeWindowsEvents:
			newGenerator = newWindowsEventsGenerator(gen, logger)
		}
		if newGenerator != nil && newGenerator.SupportsType(supportedType) {
			newGenerator.initialize()
			generators = append(generators, newGenerator)
		}
	}
	return generators
}
