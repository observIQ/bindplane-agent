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

	"github.com/observiq/observiq-otel-collector/internal/expr"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr = "sampling"

	stability = component.StabilityLevelAlpha
)

var (
	consumerCapabilities = consumer.Capabilities{MutatesData: true}
)

var once sync.Once

// NewFactory creates a new ProcessorFactory with default configuration
func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesProcessor(createTracesProcessor, stability),
		component.WithLogsProcessor(createLogsProcessor, stability),
		component.WithMetricsProcessor(createMetricsProcessor, stability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
		DropRatio:         0.5,
	}
}

func createTracesProcessor(
	ctx context.Context,
	set component.ProcessorCreateSettings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	oCfg := cfg.(*Config)

	if oCfg.Match != "" {
		return nil, fmt.Errorf("matches not supported for traces")
	}

	tmp := newSamplingProcessor(set.Logger, oCfg, nil)

	return processorhelper.NewTracesProcessor(ctx, set, cfg, nextConsumer, tmp.processTraces, processorhelper.WithCapabilities(consumerCapabilities))
}

func createLogsProcessor(
	ctx context.Context,
	set component.ProcessorCreateSettings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (component.LogsProcessor, error) {
	oCfg := cfg.(*Config)

	var match *expr.Expression
	var err error
	if oCfg.Match != "" {
		match, err = expr.CreateBoolExpression(oCfg.Match)
		if err != nil {
			return nil, fmt.Errorf("invalid match expression: %w", err)
		}
	}

	tmp := newSamplingProcessor(set.Logger, oCfg, match)

	return processorhelper.NewLogsProcessor(ctx, set, cfg, nextConsumer, tmp.processLogs, processorhelper.WithCapabilities(consumerCapabilities))
}

func createMetricsProcessor(
	ctx context.Context,
	set component.ProcessorCreateSettings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (component.MetricsProcessor, error) {
	oCfg := cfg.(*Config)

	var match *expr.Expression
	var err error
	if oCfg.Match != "" {
		match, err = expr.CreateBoolExpression(oCfg.Match)
		if err != nil {
			return nil, fmt.Errorf("invalid match expression: %w", err)
		}
	}
	tmp := newSamplingProcessor(set.Logger, oCfg, match)

	return processorhelper.NewMetricsProcessor(ctx, set, cfg, nextConsumer, tmp.processMetrics, processorhelper.WithCapabilities(consumerCapabilities))
}
