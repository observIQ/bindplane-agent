package chronicleexporter

import (
	"context"
	"errors"

	"github.com/observiq/bindplane-agent/exporter/chronicleexporter/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// TypeStr for the exporter type.
const TypeStr = "chronicle"

// Factory is the factory for the Chronicle exporter.
type Factory struct {
}

// NewFactory creates a new Chronicle exporter factory.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		TypeStr,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.LogsStability))
}

// createDefaultConfig creates the default configuration for the exporter.
func createDefaultConfig() component.Config {
	return &Config{
		TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
		QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
	}
}

// createLogsExporter creates a new log exporter based on this config.
func createLogsExporter(
	ctx context.Context,
	params exporter.CreateSettings,
	cfg component.Config,
) (exporter.Logs, error) {
	chronicleCfg, ok := cfg.(*Config)
	if !ok {
		return nil, consumererror.NewPermanent(errors.New("invalid config type"))
	}

	// Validate the config
	if err := chronicleCfg.Validate(); err != nil {
		return nil, err
	}

	exp, err := newExporter(chronicleCfg, params)
	if err != nil {
		return nil, err
	}

	return exporterhelper.NewLogsExporter(
		ctx,
		params,
		chronicleCfg,
		exp.logsDataPusher,
		exporterhelper.WithCapabilities(exp.Capabilities()),
		exporterhelper.WithTimeout(chronicleCfg.TimeoutSettings),
		exporterhelper.WithQueue(chronicleCfg.QueueSettings),
		exporterhelper.WithRetry(chronicleCfg.RetrySettings),
	)
}
