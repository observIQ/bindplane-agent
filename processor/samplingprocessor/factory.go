// Copyright  observIQ, Inc.
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

package samplingprocessor

import (
	"context"
	"fmt"
	"sync"

	"github.com/observiq/bindplane-otel-collector/expr"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

// componentType is the value of the "type" key in configuration.
var componentType = component.MustNewType("sampling")

const (
	stability = component.StabilityLevelAlpha
)

var (
	consumerCapabilities = consumer.Capabilities{MutatesData: true}
)

var once sync.Once

// NewFactory creates a new ProcessorFactory with default configuration
func NewFactory() processor.Factory {
	return processor.NewFactory(
		componentType,
		createDefaultConfig,
		processor.WithTraces(createTracesProcessor, stability),
		processor.WithLogs(createLogsProcessor, stability),
		processor.WithMetrics(createMetricsProcessor, stability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		DropRatio: 0.5,
		Condition: "true",
	}
}

func createTracesProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (processor.Traces, error) {
	oCfg := cfg.(*Config)
	condition, err := expr.NewOTTLSpanCondition(oCfg.Condition, set.TelemetrySettings)
	if err != nil {
		return nil, fmt.Errorf("invalid condition: %w", err)
	}
	sp := newTracesSamplingProcessor(set.Logger, oCfg, condition)

	return processorhelper.NewTraces(ctx, set, cfg, nextConsumer, sp.processTraces, processorhelper.WithCapabilities(consumerCapabilities))
}

func createLogsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (processor.Logs, error) {
	oCfg := cfg.(*Config)
	condition, err := expr.NewOTTLLogRecordCondition(oCfg.Condition, set.TelemetrySettings)
	if err != nil {
		return nil, fmt.Errorf("invalid condition: %w", err)
	}
	tmp := newLogsSamplingProcessor(set.Logger, oCfg, condition)

	return processorhelper.NewLogs(ctx, set, cfg, nextConsumer, tmp.processLogs, processorhelper.WithCapabilities(consumerCapabilities))
}

func createMetricsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (processor.Metrics, error) {
	oCfg := cfg.(*Config)
	condition, err := expr.NewOTTLMetricCondition(oCfg.Condition, set.TelemetrySettings)
	if err != nil {
		return nil, fmt.Errorf("invalid condition: %w", err)
	}
	tmp := newMetricsSamplingProcessor(set.Logger, oCfg, condition)

	return processorhelper.NewMetrics(ctx, set, cfg, nextConsumer, tmp.processMetrics, processorhelper.WithCapabilities(consumerCapabilities))
}
