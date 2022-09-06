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

// Package report contains reporters for collecting specific information about the collector-
package report

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/observiq/observiq-otel-collector/internal/report/snapshot"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// snapShotKind is the kind for the snapshot reporter
var snapShotKind = "snapshot"

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
var _ snapshot.Snapshotter = (*SnapshotReporter)(nil)

// SnapshotReporter tracks and reports snapshots
type SnapshotReporter struct {
	client Client

	// minPayloadSize is the minimum number of items to be in a payload
	minPayloadSize int

	// Buffers
	logBuffers    map[string]*snapshot.LogBuffer
	metricBuffers map[string]*snapshot.MetricBuffer
	traceBuffers  map[string]*snapshot.TraceBuffer
}

// NewSnapshotReporter creates a new SnapshotReporter with the associated client
func NewSnapshotReporter(client Client) *SnapshotReporter {
	return &SnapshotReporter{
		client:         client,
		minPayloadSize: 100,
		logBuffers:     make(map[string]*snapshot.LogBuffer),
		metricBuffers:  make(map[string]*snapshot.MetricBuffer),
		traceBuffers:   make(map[string]*snapshot.TraceBuffer),
	}
}

// Kind returns kind of the reporter
func (s *SnapshotReporter) Kind() string {
	return snapShotKind
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
	buffer, ok := s.logBuffers[componentID]
	if !ok {
		buffer = snapshot.NewLogBuffer(s.minPayloadSize)
	}

	buffer.Add(ld)
	s.logBuffers[componentID] = buffer
}

// SaveTraces saves off traces in a snapshot to be reported later
func (s *SnapshotReporter) SaveTraces(componentID string, td ptrace.Traces) {
	buffer, ok := s.traceBuffers[componentID]
	if !ok {
		buffer = snapshot.NewTraceBuffer(s.minPayloadSize)
	}

	buffer.Add(td)
	s.traceBuffers[componentID] = buffer
}

// SaveMetrics saves off metrics in a snapshot to be reported later
func (s *SnapshotReporter) SaveMetrics(componentID string, md pmetric.Metrics) {
	buffer, ok := s.metricBuffers[componentID]
	if !ok {
		buffer = snapshot.NewMetricBuffer(s.minPayloadSize)
	}

	buffer.Add(md)
	s.metricBuffers[componentID] = buffer
}

// prepRequestPayload based on the pipelineType will return a marshaled proto of the OTLP data types for the componentID
func (s *SnapshotReporter) prepRequestPayload(componentID, pipelineType string) (payload []byte, err error) {
	switch pipelineType {
	case "logs":
		buffer, ok := s.logBuffers[componentID]
		if !ok {
			return []byte{}, nil
		}

		payload, err = buffer.ConstructPayload()
	case "metrics":
		buffer, ok := s.metricBuffers[componentID]
		if !ok {
			return []byte{}, nil
		}

		payload, err = buffer.ConstructPayload()
	case "traces":
		buffer, ok := s.traceBuffers[componentID]
		if !ok {
			return []byte{}, nil
		}

		payload, err = buffer.ConstructPayload()
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
