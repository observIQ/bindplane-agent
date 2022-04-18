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

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/model/pdata"
)

// Exporter is a google cloud exporter wrapped with additional functionality
type Exporter struct {
	metricsProcessors []component.MetricsProcessor
	metricsExporter   component.MetricsExporter
	metricsConsumer   consumer.Metrics

	logsProcessors []component.LogsProcessor
	logsExporter   component.LogsExporter
	logsConsumer   consumer.Logs

	tracesProcessors []component.TracesProcessor
	tracesExporter   component.TracesExporter
	tracesConsumer   consumer.Traces
}

// ConsumeMetrics consumes metrics
func (e *Exporter) ConsumeMetrics(ctx context.Context, md pdata.Metrics) error {
	if e.metricsConsumer == nil {
		return nil
	}
	return e.metricsConsumer.ConsumeMetrics(ctx, md)
}

// ConsumeTraces consumes traces
func (e *Exporter) ConsumeTraces(ctx context.Context, td pdata.Traces) error {
	if e.tracesConsumer == nil {
		return nil
	}
	return e.tracesConsumer.ConsumeTraces(ctx, td)
}

// ConsumeLogs consumes logs
func (e *Exporter) ConsumeLogs(ctx context.Context, ld pdata.Logs) error {
	if e.logsConsumer == nil {
		return nil
	}
	return e.logsConsumer.ConsumeLogs(ctx, ld)
}

// Capabilities returns the capabilities of the exporter
func (e *Exporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Start starts the exporter
func (e *Exporter) Start(ctx context.Context, host component.Host) error {
	if e.tracesExporter != nil {
		if err := e.tracesExporter.Start(ctx, host); err != nil {
			return fmt.Errorf("failed to start traces exporter: %w", err)
		}
	}

	if e.logsExporter != nil {
		if err := e.logsExporter.Start(ctx, host); err != nil {
			return fmt.Errorf("failed to start logs exporter: %w", err)
		}
	}

	if e.metricsExporter != nil {
		if err := e.metricsExporter.Start(ctx, host); err != nil {
			return fmt.Errorf("failed to start metrics exporter: %w", err)
		}
	}

	for _, processor := range e.tracesProcessors {
		if err := processor.Start(ctx, host); err != nil {
			return fmt.Errorf("failed to start traces processor: %w", err)
		}
	}

	for _, processor := range e.logsProcessors {
		if err := processor.Start(ctx, host); err != nil {
			return fmt.Errorf("failed to start logs processor: %w", err)
		}
	}

	for _, processor := range e.metricsProcessors {
		if err := processor.Start(ctx, host); err != nil {
			return fmt.Errorf("failed to start metrics processor: %w", err)
		}
	}

	return nil
}

// Shutdown will shutdown the exporter
func (e *Exporter) Shutdown(ctx context.Context) error {
	for i := len(e.tracesProcessors) - 1; i >= 0; i-- {
		if err := e.tracesProcessors[i].Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown traces processor: %w", err)
		}
	}

	for i := len(e.logsProcessors) - 1; i >= 0; i-- {
		if err := e.logsProcessors[i].Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown logs processor: %w", err)
		}
	}

	for i := len(e.metricsProcessors) - 1; i >= 0; i-- {
		if err := e.metricsProcessors[i].Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown metrics processor: %w", err)
		}
	}

	if e.tracesExporter != nil {
		if err := e.tracesExporter.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown traces exporter: %w", err)
		}
	}

	if e.logsExporter != nil {
		if err := e.logsExporter.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown logs exporter: %w", err)
		}
	}

	if e.metricsExporter != nil {
		if err := e.metricsExporter.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown metrics exporter: %w", err)
		}
	}

	return nil
}
