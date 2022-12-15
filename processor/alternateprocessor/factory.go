package alternateprocessor

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

const (
	typeStr = "alternate"
)

var (
	errInvalidConfig = errors.New("invalid configuration supplied")
)

type alternateProcessorFactory struct{}

// NewFactory returns a factory that will create the alternate processor
func NewFactory() component.ProcessorFactory {
	apf := &alternateProcessorFactory{}
	return component.NewProcessorFactory(
		typeStr,
		createDefaultConfig,
		component.WithLogsProcessor(apf.createLogsProcessor, component.StabilityLevelAlpha),
		component.WithMetricsProcessor(apf.createMetricsProcessor, component.StabilityLevelAlpha),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
		Logs: &AlternateRoute{
			Enabled:             false,
			AggregationInterval: defaultAggregationInterval,
		},
		Metrics: &AlternateRoute{
			Enabled:             false,
			AggregationInterval: defaultAggregationInterval,
		},
		Traces: &AlternateRoute{
			Enabled:             false,
			AggregationInterval: defaultAggregationInterval,
		},
	}
}

func (apf *alternateProcessorFactory) createLogsProcessor(ctx context.Context, params component.ProcessorCreateSettings, cfg component.Config, consumer consumer.Logs) (component.LogsProcessor, error) {
	pConf, ok := cfg.(*Config)
	if !ok {
		return nil, errInvalidConfig
	}
	return newProcessor(pConf, params.Logger, withLogsConsumer(consumer))
}

func (apf *alternateProcessorFactory) createMetricsProcessor(_ context.Context, params component.ProcessorCreateSettings, cfg component.Config, consumer consumer.Metrics) (component.MetricsProcessor, error) {
	pConf, ok := cfg.(*Config)
	if !ok {
		return nil, errInvalidConfig
	}
	return newProcessor(pConf, params.Logger, withMetricsConsumer(consumer))
}

func (apf *alternateProcessorFactory) createTraceProcessor(ctx context.Context, params component.ProcessorCreateSettings, cfg component.Config, consumer consumer.Traces) (component.TracesProcessor, error) {
	pConf, ok := cfg.(*Config)
	if !ok {
		return nil, errInvalidConfig
	}
	return newProcessor(pConf, params.Logger, withTracesProcessor(consumer))
}
