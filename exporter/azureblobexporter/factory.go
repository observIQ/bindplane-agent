package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"context"
	"errors"

	"github.com/observiq/bindplane-agent/exporter/azureblobexporter/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// NewFactory creates a factory for Azure Blob Exporter
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithMetrics(createMetricsExporter, metadata.MetricsStability),
		exporter.WithLogs(createLogsExporter, metadata.LogsStability),
		exporter.WithTraces(createTracesExporter, metadata.TracesStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
		QueueSettings:   exporterhelper.NewDefaultQueueSettings(),
		RetrySettings:   exporterhelper.NewDefaultRetrySettings(),
		Partition:       minute,
	}
}

func createMetricsExporter(ctx context.Context, params exporter.CreateSettings, config component.Config) (exporter.Metrics, error) {
	_, ok := config.(*Config)
	if !ok {
		return nil, errors.New("not an Azure Blob config")
	}
	return nil, nil
}

func createLogsExporter(ctx context.Context, params exporter.CreateSettings, config component.Config) (exporter.Logs, error) {
	_, ok := config.(*Config)
	if !ok {
		return nil, errors.New("not an Azure Blob config")
	}
	return nil, nil
}

func createTracesExporter(ctx context.Context, params exporter.CreateSettings, config component.Config) (exporter.Traces, error) {
	_, ok := config.(*Config)
	if !ok {
		return nil, errors.New("not an Azure Blob config")
	}
	return nil, nil
}
