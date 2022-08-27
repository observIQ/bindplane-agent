package snapshotprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr = "snapshotprocessor"

	stability = component.StabilityLevelAlpha
)

var consumerCapabilities = consumer.Capabilities{MutatesData: false}

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

func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
		Enabled:           true,
	}
}

func createTracesProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	oCfg := cfg.(*Config)
	sp, err := newSnapshotProcessor(set.Logger, oCfg, cfg.ID().String())
	if err != nil {
		return nil, err
	}

	return processorhelper.NewTracesProcessor(cfg, nextConsumer, sp.processTraces, processorhelper.WithCapabilities(consumerCapabilities))
}

func createLogsProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Logs,
) (component.LogsProcessor, error) {
	oCfg := cfg.(*Config)
	sp, err := newSnapshotProcessor(set.Logger, oCfg, cfg.ID().String())
	if err != nil {
		return nil, err
	}

	return processorhelper.NewLogsProcessor(cfg, nextConsumer, sp.processLogs, processorhelper.WithCapabilities(consumerCapabilities))
}

func createMetricsProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Metrics,
) (component.MetricsProcessor, error) {
	oCfg := cfg.(*Config)
	sp, err := newSnapshotProcessor(set.Logger, oCfg, cfg.ID().String())
	if err != nil {
		return nil, err
	}

	return processorhelper.NewMetricsProcessor(cfg, nextConsumer, sp.processMetrics, processorhelper.WithCapabilities(consumerCapabilities))
}
