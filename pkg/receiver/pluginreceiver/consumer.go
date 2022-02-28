package pluginreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Consumer is the consumer placed at the end of a pipeline
type Consumer struct {
	consumer.Logs
	consumer.Metrics
	consumer.Traces
}

// createFactory creates a factory that will wrap this consumer as an exporter
func (c *Consumer) createFactory() component.ExporterFactory {
	createDefaultConfig := func() config.Exporter {
		return &struct{ config.ExporterSettings }{
			ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		}
	}

	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
		component.WithLogsExporter(c.createLogsExporter),
		component.WithTracesExporter(c.createTracesExporter),
		component.WithMetricsExporter(c.createMetricsExporter))
}

// createLogsExporter creates a logs exporter using this consumer
func (c *Consumer) createLogsExporter(_ context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.LogsExporter, error) {
	return exporterhelper.NewLogsExporter(cfg, set, c.Logs.ConsumeLogs)
}

// createTracesExporter creates a traces exporter using this consumer
func (c *Consumer) createTracesExporter(_ context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.TracesExporter, error) {
	return exporterhelper.NewTracesExporter(cfg, set, c.Traces.ConsumeTraces)
}

// createMetricsExporter creates a metrics exporter using this consumer
func (c *Consumer) createMetricsExporter(_ context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.MetricsExporter, error) {
	return exporterhelper.NewMetricsExporter(cfg, set, c.Metrics.ConsumeMetrics)
}
