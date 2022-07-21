package bindplaneexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	typeStr   = "bindplane"
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a new bindplane exporter factory
func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
	)
}

// createMetricsExporter creates a metrics exporter based on this config.
func createMetricsExporter(_ context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.MetricsExporter, error) {
	eCfg := cfg.(*Config)
	exporter := Exporter{}
	return exporterhelper.NewMetricsExporter(
		cfg,
		set,
		exporter.consumeMetrics,
		exporterhelper.WithTimeout(eCfg.TimeoutSettings),
		exporterhelper.WithQueue(eCfg.QueueSettings),
		exporterhelper.WithRetry(eCfg.RetrySettings),
	)
}

// createLogExporter creates a logs exporter based on this config.
func createLogsExporter(_ context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.LogsExporter, error) {
	eCfg := cfg.(*Config)
	exporter := Exporter{}
	return exporterhelper.NewLogsExporter(
		cfg,
		set,
		exporter.consumeLogs,
		exporterhelper.WithTimeout(eCfg.TimeoutSettings),
		exporterhelper.WithQueue(eCfg.QueueSettings),
		exporterhelper.WithRetry(eCfg.RetrySettings),
	)
}

// createTracesExporter creates a traces exporter based on this config.
func createTracesExporter(_ context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.TracesExporter, error) {
	eCfg := cfg.(*Config)
	exporter := Exporter{}
	return exporterhelper.NewTracesExporter(
		cfg,
		set,
		exporter.consumeTraces,
		exporterhelper.WithTimeout(eCfg.TimeoutSettings),
		exporterhelper.WithQueue(eCfg.QueueSettings),
		exporterhelper.WithRetry(eCfg.RetrySettings),
	)
}
