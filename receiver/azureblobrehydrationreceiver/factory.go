package azureblobrehydrationreceiver //import "github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver"

import (
	"context"
	"errors"

	"github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

// errImproperCfgType error for when an invalid config type is passed to receiver creation funcs
var errImproperCfgType = errors.New("improper config type")

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability),
		receiver.WithLogs(createLogsReceiver, metadata.LogsStability),
		receiver.WithTraces(createTracesReceiver, metadata.TracesStability),
	)
}

// createDefaultConfig creates a default configuration
func createDefaultConfig() component.Config {
	return &Config{
		DeleteOnRead: false,
	}
}

// createMetricsReceiver creates a metrics receiver
func createMetricsReceiver(_ context.Context, params receiver.CreateSettings, conf component.Config, con consumer.Metrics) (receiver.Metrics, error) {
	cfg, ok := conf.(*Config)
	if !ok {
		return nil, errImproperCfgType
	}

	return newMetricsReceiver(params.Logger, cfg, con)
}

// createLogsReceiver creates a logs receiver
func createLogsReceiver(_ context.Context, params receiver.CreateSettings, conf component.Config, con consumer.Logs) (receiver.Logs, error) {
	cfg, ok := conf.(*Config)
	if !ok {
		return nil, errImproperCfgType
	}

	return newLogsReceiver(params.Logger, cfg, con)
}

// createTracesReceiver creates a traces receiver
func createTracesReceiver(_ context.Context, params receiver.CreateSettings, conf component.Config, con consumer.Traces) (receiver.Traces, error) {
	cfg, ok := conf.(*Config)
	if !ok {
		return nil, errImproperCfgType
	}

	return newTracesReceiver(params.Logger, cfg, con)
}
