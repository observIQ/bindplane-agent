package snapshotprocessor

import (
	"context"

	"github.com/observiq/observiq-otel-collector/internal/report"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// getSnapshotReporter is function for retrieving the SnapshotReporter.
// Meant to be overridden for tests.
var getSnapshotReporter func() *report.SnapshotReporter = report.GetSnapshotReporter

type snapshotProcessor struct {
	logger           *zap.Logger
	enabled          bool
	snapShotReporter *report.SnapshotReporter
	processorID      string
}

func newSnapshotProcessor(logger *zap.Logger, cfg *Config, processorID string) (*snapshotProcessor, error) {
	return &snapshotProcessor{
		logger:           logger,
		enabled:          cfg.Enabled,
		snapShotReporter: getSnapshotReporter(),
		processorID:      processorID,
	}, nil
}

func (sp *snapshotProcessor) processTraces(_ context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	if sp.enabled {
		sp.snapShotReporter.SaveTraces(sp.processorID, td.Clone())
	}

	return td, nil
}

func (sp *snapshotProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	if sp.enabled {
		sp.snapShotReporter.SaveLogs(sp.processorID, ld.Clone())
	}

	return ld, nil
}

func (sp *snapshotProcessor) processMetrics(_ context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	if sp.enabled {
		sp.snapShotReporter.SaveMetrics(sp.processorID, md.Clone())
	}

	return md, nil
}
