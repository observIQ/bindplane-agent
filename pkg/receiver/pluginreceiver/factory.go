package pluginreceiver

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

const (
	typeStr = "plugin"
)

// NewFactory creates a new factory for the plugin receiver.
func NewFactory() component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesReceiver(createTracesReceiver),
		component.WithMetricsReceiver(createMetricsReceiver),
		component.WithLogsReceiver(createLogReceiver))
}

// createDefaultConfig creates the default configuration for the plugin receiver.
func createDefaultConfig() config.Receiver {
	return &Config{
		ReceiverSettings: config.NewReceiverSettings(config.NewComponentID(typeStr)),
	}
}

// CreateTracesReceiver creates a trace receiver from the supplied parameters
func createTracesReceiver(_ context.Context, set component.ReceiverCreateSettings, cfg config.Receiver, traceConsumer consumer.Traces) (component.TracesReceiver, error) {
	consumer := Consumer{Traces: traceConsumer}
	return createReceiver(set, cfg, consumer)
}

// CreateMetricsReceiver creates a metrics receiver from the supplied parameters
func createMetricsReceiver(_ context.Context, set component.ReceiverCreateSettings, cfg config.Receiver, metricConsumer consumer.Metrics) (component.MetricsReceiver, error) {
	consumer := Consumer{Metrics: metricConsumer}
	return createReceiver(set, cfg, consumer)
}

// CreateLogReceiver creates a log receiver from the supplied parameters
func createLogReceiver(_ context.Context, set component.ReceiverCreateSettings, cfg config.Receiver, logConsumer consumer.Logs) (component.LogsReceiver, error) {
	consumer := Consumer{Logs: logConsumer}
	return createReceiver(set, cfg, consumer)
}

// createReceiver creates a receiver from the supplied parameters
func createReceiver(set component.ReceiverCreateSettings, cfg config.Receiver, consumer Consumer) (*Receiver, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("invalid config type")
	}

	return &Receiver{
		template:   config.Template,
		parameters: config.Parameters,
		consumer:   consumer,
		telemetry:  set.TelemetrySettings,
		buildInfo:  set.BuildInfo,
	}, nil
}
