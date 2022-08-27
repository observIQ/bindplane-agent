package snapshotprocessor

import (
	"context"
	"fmt"

	"github.com/observiq/observiq-otel-collector/internal/report"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type snapshotProcessor struct {
	logger           *zap.Logger
	enabled          bool
	snapShotReporter *report.SnapshotReporter
	processorID      string
}

func newSnapshotProcessor(logger *zap.Logger, cfg *Config, processorID string) (*snapshotProcessor, error) {
	reporter, err := report.GetSnapshotReporter()
	if err != nil {
		return nil, fmt.Errorf("error retrieving SnapshotReporter: %w", err)
	}

	return &snapshotProcessor{
		logger:           logger,
		enabled:          cfg.Enabled,
		snapShotReporter: reporter,
		processorID:      processorID,
	}, nil
}

func (sp *snapshotProcessor) processTraces(_ context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	if sp.enabled {
		sp.snapShotReporter.ReportTraces(sp.processorID, td)
	}

	return td, nil
}

func (sp *snapshotProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	if sp.enabled {
		sp.snapShotReporter.ReportLogs(sp.processorID, ld)
	}

	return ld, nil
}

func (sp *snapshotProcessor) processMetrics(_ context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	if sp.enabled {
		sp.snapShotReporter.ReportMetrics(sp.processorID, md)
	}

	return md, nil
}
