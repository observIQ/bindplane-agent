package splunksearchapireceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

var (
	typeStr = component.MustNewType("splunksearchapi")
)

func createDefaultConfig() component.Config {
	return &Config{}
}

func createLogsReceiver(_ context.Context,
	params receiver.Settings,
	cfg component.Config,
	consumer consumer.Logs,
) (receiver.Logs, error) {
	logger := params.Logger
	ssapirConfig := cfg.(*Config)
	ssapir := &splunksearchapireceiver{
		logger:       logger,
		logsConsumer: consumer,
		config:       ssapirConfig,
	}
	return ssapir, nil
}

func NewFactory() receiver.Factory {
	return receiver.NewFactory(typeStr, createDefaultConfig, receiver.WithLogs(createLogsReceiver, component.StabilityLevelAlpha))
}
