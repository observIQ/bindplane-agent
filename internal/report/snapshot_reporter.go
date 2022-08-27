// Package report contains reporters for collecting specific information about the collector
package report

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// snapShotType is the reporterType for the snapshot reporter
var snapShotType ReporterType = "snapshot"

var _ Reporter = (*SnapshotReporter)(nil)

// SnapshotReporter tracks and reports snapshots
type SnapshotReporter struct {
	enabled bool

	// minPayloadSize is the minimum number of items to be in a payload
	minPayloadSize int

	// Buffers
	logBuffers    map[string][]plog.Logs
	metricBuffers map[string][]pmetric.Metrics
	traceBuffers  map[string][]ptrace.Traces
}

// NewSnapshotReporter creates a new SnapshotReporter with the associated client
func NewSnapshotReporter(minPayloadSize int, enabled bool) *SnapshotReporter {
	return &SnapshotReporter{
		// TODO fill out
		enabled:        enabled,
		minPayloadSize: minPayloadSize,
		logBuffers:     make(map[string][]plog.Logs),
		metricBuffers:  make(map[string][]pmetric.Metrics),
		traceBuffers:   make(map[string][]ptrace.Traces),
	}
}

// Type returns type of the reporter
func (s *SnapshotReporter) Type() ReporterType {
	return snapShotType
}

// Start kicks off reporting snapshots via the client
func (s *SnapshotReporter) Start(_ Client) error {
	// TODO send data
	return nil
}

// Stop stops reporting snapshots and clears all buffers
func (s *SnapshotReporter) Stop(context.Context) error {
	for k := range s.logBuffers {
		delete(s.logBuffers, k)
	}
	for k := range s.traceBuffers {
		delete(s.logBuffers, k)
	}
	for k := range s.metricBuffers {
		delete(s.logBuffers, k)
	}

	// TODO
	return nil
}

// IsEnabled implements Reporter
func (s *SnapshotReporter) IsEnabled() bool {
	return s.enabled
}

// ReportLogs reports logs to be sent to platform
func (s *SnapshotReporter) ReportLogs(componentID string, ld plog.Logs) {
	componentLogs, ok := s.logBuffers[componentID]
	if !ok {
		componentLogs = make([]plog.Logs, 0)
	}

	currentPayloadSize := 0
	for _, logs := range componentLogs {
		currentPayloadSize += logs.LogRecordCount()
	}

	componentLogs = insertPayload(componentLogs, ld, ld.LogRecordCount(), currentPayloadSize, s.minPayloadSize)

	s.logBuffers[componentID] = componentLogs
}

// ReportTraces reports traces to be sent to platform
func (s *SnapshotReporter) ReportTraces(componentID string, td ptrace.Traces) {
	componentTraces, ok := s.traceBuffers[componentID]
	if !ok {
		componentTraces = make([]ptrace.Traces, 0)
	}

	currentPayloadSize := 0
	for _, traces := range componentTraces {
		currentPayloadSize += traces.SpanCount()
	}

	componentTraces = insertPayload(componentTraces, td, td.SpanCount(), currentPayloadSize, s.minPayloadSize)

	s.traceBuffers[componentID] = componentTraces
}

// ReportMetrics reports metrics to be sent to platform
func (s *SnapshotReporter) ReportMetrics(componentID string, md pmetric.Metrics) {
	componentMetrics, ok := s.metricBuffers[componentID]
	if !ok {
		componentMetrics = make([]pmetric.Metrics, 0)
	}

	currentPayloadSize := 0
	for _, metrics := range componentMetrics {
		currentPayloadSize += metrics.DataPointCount()
	}

	componentMetrics = insertPayload(componentMetrics, md, md.DataPointCount(), currentPayloadSize, s.minPayloadSize)

	s.metricBuffers[componentID] = componentMetrics
}

func insertPayload[T plog.Logs | pmetric.Metrics | ptrace.Traces](buffer []T, payload T, payloadSize, currentPayloadSize, minPayloadSize int) []T {
	switch {
	// The number of payload items is more than minPayloadSize so reset this to just this Log set
	case payloadSize > minPayloadSize:
		buffer = []T{payload}

	// Haven't reached full size yet so add this
	case payloadSize+currentPayloadSize < minPayloadSize:
		buffer = append(buffer, payload)

	// Adding this will put us over minimum so remove the oldest record and add the new one
	case payloadSize+currentPayloadSize >= minPayloadSize:
		buffer = append(buffer, payload)
		buffer = buffer[1:]
	}

	return buffer
}
