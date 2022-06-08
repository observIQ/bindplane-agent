// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package googlecloudexporter

import (
	"context"
	"fmt"

	gcp "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
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
	exporterConfig.setClientOptions()

	gcpExporter, err := gcpFactory.CreateMetricsExporter(ctx, set, exporterConfig.GCPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	processorSettings := component.ProcessorCreateSettings{
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	exporterConfig.BatchConfig.SetIDName(exporterConfig.ID().String())
	factory := batchprocessor.NewFactory()
	batchProcessor, err := factory.CreateMetricsProcessor(ctx, processorSettings, exporterConfig.BatchConfig, gcpExporter)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch processor: %w", err)
	}

	return &exporter{
		metricsBatcher:  batchProcessor,
		metricsExporter: gcpExporter,
	}, nil
}

// createLogExporter creates a logs exporter based on this config.
func createLogsExporter(ctx context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.LogsExporter, error) {
	exporterConfig := cfg.(*Config)
	exporterConfig.setClientOptions()

	gcpExporter, err := gcpFactory.CreateLogsExporter(ctx, set, exporterConfig.GCPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs exporter: %w", err)
	}

	processorSettings := component.ProcessorCreateSettings{
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	exporterConfig.BatchConfig.SetIDName(exporterConfig.ID().String())
	factory := batchprocessor.NewFactory()
	batchProcessor, err := factory.CreateLogsProcessor(ctx, processorSettings, exporterConfig.BatchConfig, gcpExporter)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch processor: %w", err)
	}

	return &exporter{
		logsBatcher:  batchProcessor,
		logsExporter: gcpExporter,
	}, nil
}

// createTracesExporter creates a traces exporter based on this config.
func createTracesExporter(ctx context.Context, set component.ExporterCreateSettings, cfg config.Exporter) (component.TracesExporter, error) {
	exporterConfig := cfg.(*Config)
	exporterConfig.setClientOptions()

	gcpExporter, err := gcpFactory.CreateTracesExporter(ctx, set, exporterConfig.GCPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create traces exporter: %w", err)
	}

	processorSettings := component.ProcessorCreateSettings{
		TelemetrySettings: set.TelemetrySettings,
		BuildInfo:         set.BuildInfo,
	}

	exporterConfig.BatchConfig.SetIDName(exporterConfig.ID().String())
	factory := batchprocessor.NewFactory()
	batchProcessor, err := factory.CreateTracesProcessor(ctx, processorSettings, exporterConfig.BatchConfig, gcpExporter)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch processor: %w", err)
	}

	return &exporter{
		tracesBatcher:  batchProcessor,
		tracesExporter: gcpExporter,
	}, nil
}
