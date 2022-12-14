package routerprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
)

const (
	// typeStr is the value of the "type" key in the configuration.
	typeStr = "router"

	// stability is the current state of the processor.
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a new factory for the processor.
func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig,
		component.WithLogsProcessor(createLogsProcessor, stability),
	)
}

// createLogsProcessor creates a log processor.
func createLogsProcessor(_ context.Context, params component.ProcessorCreateSettings, cfg component.Config, consumer consumer.Logs) (component.LogsProcessor, error) {
	processorCfg, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type: %+v", cfg)
	}

	return newProcessor(processorCfg, consumer, params.Logger)
}
