// Package report contains reporters for collecting specific information about the collector-
package report

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// snapShotType is the reporterType for the snapshot reporter
var snapShotType ReporterKind = "snapshot"

// snapshotConfig specifies what snapshots to collect
type snapshotConfig struct {
	// Endpoint is where to send the snapshots
	Endpoint *endpointConfig `yaml:"endpoint"`

	// Processor is the full ComponentID of the snapshot processor
	Processor string `yaml:"processor"`

	// PipelineType will be "logs", "metrics", or "traces"
	PipelineType string `yaml:"pipeline_type"`
}

// endpointConfig is the configuration of a specific endpoint and full headers to include
type endpointConfig struct {
	URL     string      `yaml:"url"`
	Headers http.Header `yaml:"headers"`
}

var _ Reporter = (*SnapshotReporter)(nil)

// SnapshotReporter tracks and reports snapshots
type SnapshotReporter struct {
	client Client

	// minPayloadSize is the minimum number of items to be in a payload
	minPayloadSize int

	// Buffers
	logBuffers    map[string][]plog.Logs
	metricBuffers map[string][]pmetric.Metrics
	traceBuffers  map[string][]ptrace.Traces
}

// NewSnapshotReporter creates a new SnapshotReporter with the associated client
func NewSnapshotReporter(client Client) *SnapshotReporter {
	return &SnapshotReporter{
		client:         client,
		minPayloadSize: 100,
		logBuffers:     make(map[string][]plog.Logs),
		metricBuffers:  make(map[string][]pmetric.Metrics),
		traceBuffers:   make(map[string][]ptrace.Traces),
	}
}

// Type returns type of the reporter
func (s *SnapshotReporter) Type() ReporterKind {
	return snapShotType
}

// Report applies the new configuration and reports snapshots specified in it
func (s *SnapshotReporter) Report(cfg any) error {
	ssCfg, ok := cfg.(*snapshotConfig)
	if !ok {
		return errors.New("invalid config type")
	}

	// Gather payload
	payload, err := s.prepRequestPayload(ssCfg.Processor, ssCfg.PipelineType)

	// Compress
	compressedPayload, err := compress(payload)
	if err != nil {
		return fmt.Errorf("failed to compress payload: %w", err)
	}

	// Prep request
	req, err := http.NewRequest(http.MethodPost, ssCfg.Endpoint.URL, bytes.NewReader(compressedPayload))
	if err != nil {
		return fmt.Errorf("failed to construct snapshot request: %w", err)
	}
	// Add content headers
	req.Header.Add("Content-Type", "application/protobuf")
	req.Header.Add("Content-Encoding", "gzip")

	// Add Component-ID header
	req.Header.Add("Component-ID", ssCfg.Processor)

	// Add headers from config
	for k, values := range ssCfg.Endpoint.Headers {
		for _, value := range values {
			req.Header.Add(k, value)
		}
	}

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("snapshot request failed: %w", err)
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf("Non 200 response for snapshot report: %d", resp.StatusCode)
	}

	return nil
}

// Stop does nothing as there is no long running process
func (s *SnapshotReporter) Stop(context.Context) error {
	return nil
}

// SaveLogs saves off logs in a snapshot to be reported later
func (s *SnapshotReporter) SaveLogs(componentID string, ld plog.Logs) {
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

// SaveTraces saves off traces in a snapshot to be reported later
func (s *SnapshotReporter) SaveTraces(componentID string, td ptrace.Traces) {
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

// SaveMetrics saves off metrics in a snapshot to be reported later
func (s *SnapshotReporter) SaveMetrics(componentID string, md pmetric.Metrics) {
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

// prepRequestPayload based on the pipelineType will return a marshaled proto of the OTLP data types for the componentID
func (s *SnapshotReporter) prepRequestPayload(componentID, pipelineType string) (payload []byte, err error) {
	switch pipelineType {
	case "logs":
		logsMarshler := plog.NewProtoMarshaler()

		payloadLogs := plog.NewLogs()
		logs := s.logBuffers[componentID]
		if len(logs) > 0 {
			// Copy logs to a single payload
			for _, ld := range logs {
				ld.ResourceLogs().CopyTo(payloadLogs.ResourceLogs())
			}
		}

		payload, err = logsMarshler.MarshalLogs(payloadLogs)
		if err != nil {
			return nil, fmt.Errorf("failed to construct payload: %w", err)
		}
	case "metrics":
		metricsMarshler := pmetric.NewProtoMarshaler()

		payloadMetrics := pmetric.NewMetrics()
		metrics := s.metricBuffers[componentID]
		if len(metrics) > 0 {
			for _, md := range metrics {
				md.ResourceMetrics().CopyTo(payloadMetrics.ResourceMetrics())
			}
		}

		payload, err = metricsMarshler.MarshalMetrics(payloadMetrics)
		if err != nil {
			return nil, fmt.Errorf("failed to construct payload: %w", err)
		}
	case "traces":
		tracesMarshler := ptrace.NewProtoMarshaler()

		payloadTraces := ptrace.NewTraces()
		traces := s.traceBuffers[componentID]
		if len(traces) > 0 {
			for _, td := range traces {
				td.ResourceSpans().CopyTo(payloadTraces.ResourceSpans())
			}
		}

		payload, err = tracesMarshler.MarshalTraces(payloadTraces)
		if err != nil {
			return nil, fmt.Errorf("failed to construct payload: %w", err)
		}
	}

	// Case where there is no snapshot data for requested component and pipeline
	if payload == nil {
		payload = make([]byte, 0)
	}

	return
}

// compress gzip compresses the data
func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
