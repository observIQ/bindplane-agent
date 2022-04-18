package googlecloudexporter

import (
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-collector/processor/normalizesumsprocessor"
	"github.com/observiq/observiq-otel-collector/processor/resourceattributetransposerprocessor"
	gcp "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourcedetectionprocessor"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor/batchprocessor"
)

// gcpFactory is the factory used to create the underlying gcp exporter
var gcpFactory = gcp.NewFactory()

// typeStr is the type of the google cloud exporter
const typeStr = "googlecloud"

// NewFactory creates a factory for the googlecloud exporter
func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
		component.WithMetricsExporter(createMetricsExporter),
		component.WithLogsExporter(createLogsExporter),
		component.WithTracesExporter(createTracesExporter),
	)
}

// createMetricsExporter creates a metrics exporter based on this config.
func createMetricsExporter(ctx context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.MetricsExporter, error) {
	exporterConfig := cfg.(*Config)

	exporter, err := gcpFactory.CreateMetricsExporter(ctx, set, exporterConfig.GCPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	attributerConfig := addGenericAttributes(*exporterConfig.AttributerConfig, exporterConfig.Namespace, exporterConfig.Location)
	processors := []component.MetricsProcessor{}
	processorConfigs := []config.Processor{
		exporterConfig.BatchConfig,
		exporterConfig.NormalizeConfig,
		exporterConfig.TransposerConfig,
		attributerConfig,
		exporterConfig.DetectorConfig,
	}

	processorFactories := []component.ProcessorFactory{
		batchprocessor.NewFactory(),
		normalizesumsprocessor.NewFactory(),
		resourceattributetransposerprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		resourcedetectionprocessor.NewFactory(),
	}

	processorSettings := component.ProcessorCreateSettings{
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	var consumer consumer.Metrics = exporter
	for i, processorConfig := range processorConfigs {
		processorConfig.SetIDName(exporterConfig.ID().String())
		factory := processorFactories[i]
		processor, err := factory.CreateMetricsProcessor(ctx, processorSettings, processorConfig, consumer)
		if err != nil {
			return nil, fmt.Errorf("failed to create metrics processor %s: %w", processorConfig.ID().String(), err)
		}
		processors = append(processors, processor)
		consumer = processor
	}

	return &Exporter{
		metricsProcessors: processors,
		metricsExporter:   exporter,
		metricsConsumer:   consumer,
	}, nil
}

// createLogExporter creates a logs exporter based on this config.
func createLogsExporter(ctx context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.LogsExporter, error) {
	exporterConfig := cfg.(*Config)

	exporter, err := gcpFactory.CreateLogsExporter(ctx, set, exporterConfig.GCPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs exporter: %w", err)
	}

	attributerConfig := addGenericAttributes(*exporterConfig.AttributerConfig, exporterConfig.Namespace, exporterConfig.Location)
	processors := []component.LogsProcessor{}
	processorConfigs := []config.Processor{
		exporterConfig.BatchConfig,
		attributerConfig,
		exporterConfig.DetectorConfig,
	}

	processorFactories := []component.ProcessorFactory{
		batchprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		resourcedetectionprocessor.NewFactory(),
	}

	processorSettings := component.ProcessorCreateSettings{
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	var consumer consumer.Logs = exporter
	for i, processorConfig := range processorConfigs {
		processorConfig.SetIDName(exporterConfig.ID().String())
		factory := processorFactories[i]
		processor, err := factory.CreateLogsProcessor(ctx, processorSettings, processorConfig, consumer)
		if err != nil {
			return nil, fmt.Errorf("failed to create logs processor %s: %w", processorConfig.ID().String(), err)
		}
		processors = append(processors, processor)
		consumer = processor
	}

	return &Exporter{
		logsProcessors: processors,
		logsExporter:   exporter,
		logsConsumer:   consumer,
	}, nil
}

// createTracesExporter creates a traces exporter based on this config.
func createTracesExporter(ctx context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.TracesExporter, error) {
	exporterConfig := cfg.(*Config)

	exporter, err := gcpFactory.CreateTracesExporter(ctx, set, exporterConfig.GCPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create traces exporter: %w", err)
	}

	attributerConfig := addGenericAttributes(*exporterConfig.AttributerConfig, exporterConfig.Namespace, exporterConfig.Location)
	processors := []component.TracesProcessor{}
	processorConfigs := []config.Processor{
		exporterConfig.BatchConfig,
		attributerConfig,
		exporterConfig.DetectorConfig,
	}

	processorFactories := []component.ProcessorFactory{
		batchprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		resourcedetectionprocessor.NewFactory(),
	}

	processorSettings := component.ProcessorCreateSettings{
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	var consumer consumer.Traces = exporter
	for i, processorConfig := range processorConfigs {
		processorConfig.SetIDName(exporterConfig.ID().String())
		factory := processorFactories[i]
		processor, err := factory.CreateTracesProcessor(ctx, processorSettings, processorConfig, consumer)
		if err != nil {
			return nil, fmt.Errorf("failed to create traces processor %s: %w", processorConfig.ID().String(), err)
		}
		processors = append(processors, processor)
		consumer = processor
	}

	return &Exporter{
		tracesProcessors: processors,
		tracesExporter:   exporter,
		tracesConsumer:   consumer,
	}, nil
}
