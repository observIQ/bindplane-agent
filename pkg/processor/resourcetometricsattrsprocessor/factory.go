package resourcetometricsattrsprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr = "resourcetometricsattrs"
)

// NewFactory returns a new factory for the resourcetometricsattrs processor.
func NewFactory() component.ProcessorFactory {
	return processorhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		processorhelper.WithMetrics(createMetricsProcessor),
	)
}

// createDefaultConfig returns the default config for the resourcetometricsattrs processor.
func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
	}
}

// createMetricsProcessor creates the resourcetometricsattrs processor.
func createMetricsProcessor(ctx context.Context, params component.ProcessorCreateSettings, cfg config.Processor, nextConsumer consumer.Metrics) (component.MetricsProcessor, error) {
	processorCfg, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("config was not of correct type for the processor: %+v", cfg)
	}

	return newResourceToMetricsAttributesProcessor(params.Logger, nextConsumer, processorCfg), nil
}
