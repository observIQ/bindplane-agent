package generatereceiver

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const typeStr = "generate"

// NewFactory creates a factory for the generate receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithLogs(createLogsReceiver, component.StabilityLevelAlpha),
	)
}

// createLogsReceiver creates a logs receiver.
func createLogsReceiver(
	_ context.Context,
	settings receiver.CreateSettings,
	config component.Config,
	consumer consumer.Logs,
) (receiver.Logs, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type: %#v", config)
	}

	return &Receiver{
		config:      cfg,
		logConsumer: consumer,
		logger:      settings.Logger,
	}, nil
}
