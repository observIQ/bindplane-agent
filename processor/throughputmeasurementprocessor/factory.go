package throughputmeasurementprocessor

import (
	"context"
	"sync"

	"go.opencensus.io/stats/view"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	typeStr = "throughputmeasurement"

	stability = component.StabilityLevelAlpha
)

var (
	consumerCapabilities = consumer.Capabilities{MutatesData: false}
)

var once sync.Once

func NewFactory() component.ProcessorFactory {
	once.Do(func() {
		_ = view.Register(MetricViews()...)
	})

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
		SamplingRatio:     1.0,
	}
}

func createTracesProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Traces,
) (component.TracesProcessor, error) {
	oCfg := cfg.(*Config)
	tmp := newThroughputMeasurementProcessor(set.Logger, oCfg)

	return processorhelper.NewTracesProcessor(cfg, nextConsumer, tmp.processTraces, processorhelper.WithCapabilities(consumerCapabilities))
}

func createLogsProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Logs,
) (component.LogsProcessor, error) {
	oCfg := cfg.(*Config)
	tmp := newThroughputMeasurementProcessor(set.Logger, oCfg)

	return processorhelper.NewLogsProcessor(cfg, nextConsumer, tmp.processLogs, processorhelper.WithCapabilities(consumerCapabilities))
}

func createMetricsProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Metrics,
) (component.MetricsProcessor, error) {
	oCfg := cfg.(*Config)
	tmp := newThroughputMeasurementProcessor(set.Logger, oCfg)

	return processorhelper.NewMetricsProcessor(cfg, nextConsumer, tmp.processMetrics, processorhelper.WithCapabilities(consumerCapabilities))
}
