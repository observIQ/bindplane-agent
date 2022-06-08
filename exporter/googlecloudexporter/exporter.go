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
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// exporter is a google cloud exporter wrapped with additional functionality
type exporter struct {
	metricsBatcher  component.MetricsProcessor
	metricsExporter component.MetricsExporter

	logsBatcher  component.LogsProcessor
	logsExporter component.LogsExporter

	tracesBatcher  component.TracesProcessor
	tracesExporter component.TracesExporter
}

// ConsumeMetrics consumes metrics
func (e *exporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	if e.metricsBatcher == nil {
		return nil
	}
	return e.metricsBatcher.ConsumeMetrics(ctx, md)
}

// ConsumeTraces consumes traces
func (e *exporter) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	if e.tracesBatcher == nil {
		return nil
	}
	return e.tracesBatcher.ConsumeTraces(ctx, td)
}

// ConsumeLogs consumes logs
func (e *exporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	if e.logsBatcher == nil {
		return nil
	}
	return e.logsBatcher.ConsumeLogs(ctx, ld)
}

// Capabilities returns the capabilities of the exporter
func (e *exporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Start starts the exporter
func (e *exporter) Start(ctx context.Context, host component.Host) error {
	if err := e.startTraces(ctx, host); err != nil {
		return fmt.Errorf("failed to start traces: %w", err)
	}

	if err := e.startLogs(ctx, host); err != nil {
		return fmt.Errorf("failed to start logs: %w", err)
	}

	if err := e.startMetrics(ctx, host); err != nil {
		return fmt.Errorf("failed to start metrics: %w", err)
	}

	return nil
}

// Shutdown will shutdown the exporter
func (e *exporter) Shutdown(ctx context.Context) error {
	if err := e.shutdownTraces(ctx); err != nil {
		return fmt.Errorf("failed to shutdown traces: %w", err)
	}

	if err := e.shutdownLogs(ctx); err != nil {
		return fmt.Errorf("failed to shutdown logs: %w", err)
	}

	if err := e.shutdownMetrics(ctx); err != nil {
		return fmt.Errorf("failed to shutdown metrics: %w", err)
	}

	return nil
}

// startTraces will start the exporter's trace consumer
func (e *exporter) startTraces(ctx context.Context, host component.Host) error {
	if e.tracesExporter == nil {
		return nil
	}

	if err := e.tracesExporter.Start(ctx, host); err != nil {
		return fmt.Errorf("failed to start traces exporter: %w", err)
	}

	if err := e.tracesBatcher.Start(ctx, host); err != nil {
		return fmt.Errorf("failed to start traces batcher: %w", err)
	}

	return nil
}

// shutdownTraces will shutdown the exporter's trace consumer
func (e *exporter) shutdownTraces(ctx context.Context) error {
	if e.tracesExporter == nil {
		return nil
	}

	if err := e.tracesBatcher.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown traces batcher: %w", err)
	}

	if err := e.tracesExporter.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown traces exporter: %w", err)
	}

	return nil
}

// startLogs will start the exporter's log consumer
func (e *exporter) startLogs(ctx context.Context, host component.Host) error {
	if e.logsExporter == nil {
		return nil
	}

	if err := e.logsExporter.Start(ctx, host); err != nil {
		return fmt.Errorf("failed to start logs exporter: %w", err)
	}

	if err := e.logsBatcher.Start(ctx, host); err != nil {
		return fmt.Errorf("failed to start logs batcher: %w", err)
	}

	return nil
}

// shutdownLogs will shutdown the exporter's log consumer
func (e *exporter) shutdownLogs(ctx context.Context) error {
	if e.logsExporter == nil {
		return nil
	}

	if err := e.logsBatcher.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown logs batcher: %w", err)
	}

	if err := e.logsExporter.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown logs exporter: %w", err)
	}

	return nil
}

// startMetrics will start the exporter's metric consumer
func (e *exporter) startMetrics(ctx context.Context, host component.Host) error {
	if e.metricsExporter == nil {
		return nil
	}

	if err := e.metricsExporter.Start(ctx, host); err != nil {
		return fmt.Errorf("failed to start metrics exporter: %w", err)
	}

	if err := e.metricsBatcher.Start(ctx, host); err != nil {
		return fmt.Errorf("failed to start metrics batcher: %w", err)
	}

	return nil
}

// shutdownMetrics will shutdown the exporter's metric consumer
func (e *exporter) shutdownMetrics(ctx context.Context) error {
	if e.metricsExporter == nil {
		return nil
	}

	if err := e.metricsBatcher.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown metrics batcher: %w", err)
	}

	if err := e.metricsExporter.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown metrics exporter: %w", err)
	}

	return nil
}
