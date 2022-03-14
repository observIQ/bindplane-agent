package pluginreceiver

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
)

const typeStr = "plugin"

// Config is the configuration of a plugin receiver
type Config struct {
	config.ReceiverSettings `mapstructure:",squash"`
	path                    string                 `mapstructure:"path"`
	parameters              map[string]interface{} `mapstructure:"parameters"`
}

// createDefaultConfig creates a default config for a plugin receiver
func createDefaultConfig() config.Receiver {
	return &Config{
		ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
		parameters:       make(map[string]interface{}),
	}
}

// NewFactory creates a factory for a plugin receiver
func NewFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		receiverhelper.WithLogs(createLogsReceiver),
		receiverhelper.WithMetrics(createMetricsReceiver),
		receiverhelper.WithTraces(createTracesReceiver),
	)
}

// createLogsReceiver creates a plugin receiver with a logs consumer
func createLogsReceiver(_ context.Context, set component.ReceiverCreateSettings, cfg config.Receiver, consumer consumer.Logs) (component.LogsReceiver, error) {
	emitterFactory := createLogEmitterFactory(consumer)
	return createReceiver(cfg, set, emitterFactory)
}

// createMetricsReceiver creates a plugin receiver with a metrics consumer
func createMetricsReceiver(_ context.Context, set component.ReceiverCreateSettings, cfg config.Receiver, consumer consumer.Metrics) (component.MetricsReceiver, error) {
	emitterFactory := createMetricEmitterFactory(consumer)
	return createReceiver(cfg, set, emitterFactory)
}

// createTracesReceiver creates a plugin receiver with a traces consumer
func createTracesReceiver(_ context.Context, set component.ReceiverCreateSettings, cfg config.Receiver, consumer consumer.Traces) (component.TracesReceiver, error) {
	emitterFactory := createTraceEmitterFactory(consumer)
	return createReceiver(cfg, set, emitterFactory)
}

// createReceiver creates a plugin receiver with the supplied emitter
func createReceiver(cfg config.Receiver, set component.ReceiverCreateSettings, emitterFactory component.ExporterFactory) (*Receiver, error) {
	receiverConfig, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("config is not a plugin receiver config")
	}

	plugin, err := LoadPlugin(receiverConfig.path)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	if err := plugin.CheckParameters(receiverConfig.parameters); err != nil {
		return nil, fmt.Errorf("invalid plugin parameter: %w", err)
	}

	configMap, err := plugin.RenderConfig(receiverConfig.parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to render plugin: %w", err)
	}
	configProvider := createConfigProvider(configMap)
	factories := createFactories(emitterFactory)
	factoryProvider := createFactoryProvider(factories)

	return &Receiver{
		configProvider:  configProvider,
		factoryProvider: factoryProvider,
		buildInfo:       set.BuildInfo,
		logger:          set.Logger.Named(receiverConfig.ID().String()),
		createSvc:       createDefaultService,
	}, nil
}

// createFactories creates a factories map with the emitter factory
func createFactories(emitterFactory component.ExporterFactory) component.Factories {
	return component.Factories{
		Exporters: map[config.Type]component.ExporterFactory{
			emitterFactory.Type(): emitterFactory,
		},
	}
}
