package bindplaneexporter

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// Exporter is the bindplane exporter
type Exporter struct {
}

// ConsumeMetrics will consume metrics
func (e *Exporter) ConsumeMetrics(_ context.Context, metrics pmetric.Metrics) error {
	records := getRecordsFromMetrics(metrics)
	_ = NewMetricsMessage(records, nil)
	return nil
}

// ConsumeLogs will consume logs
func (e *Exporter) ConsumeLogs(_ context.Context, logs plog.Logs) error {
	records := getRecordsFromLogs(logs)
	_ = NewLogsMessage(records, nil)
	return nil
}

// ConsumeTraces will consume traces
func (e *Exporter) ConsumeTraces(_ context.Context, traces ptrace.Traces) error {
	records := getRecordsFromTraces(traces)
	_ = NewTracesMessage(records, nil)
	return nil
}
