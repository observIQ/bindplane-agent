// Copyright observIQ, Inc.
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

package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"bytes"
	"compress/gzip"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// marshaler marshals into data bytes based on configuration
//
//go:generate mockery --name marshaler --output ./internal/mocks --with-expecter --filename mock_marshaler.go --structname MockMarshaler
type marshaler interface {
	MarshalTraces(td ptrace.Traces) ([]byte, error)
	MarshalLogs(ld plog.Logs) ([]byte, error)
	MarshalMetrics(md pmetric.Metrics) ([]byte, error)
	Format() string
}

func newMarshaler(compression compressionType) marshaler {
	base := &baseMarshaler{
		logsMarshaler:    &plog.JSONMarshaler{},
		tracesMarshaler:  &ptrace.JSONMarshaler{},
		metricsMarshaler: &pmetric.JSONMarshaler{},
	}

	switch compression {
	case gzipCompression:
		return &gzipMarshaler{
			base: base,
		}
	default:
		return base
	}
}

// baseMarshaler is the base marshaller that marshals otlp structures into JSON
type baseMarshaler struct {
	logsMarshaler    plog.Marshaler
	tracesMarshaler  ptrace.Marshaler
	metricsMarshaler pmetric.Marshaler
}

func (b *baseMarshaler) MarshalTraces(td ptrace.Traces) ([]byte, error) {
	return b.tracesMarshaler.MarshalTraces(td)
}

func (b *baseMarshaler) MarshalLogs(ld plog.Logs) ([]byte, error) {
	return b.logsMarshaler.MarshalLogs(ld)
}

func (b *baseMarshaler) MarshalMetrics(md pmetric.Metrics) ([]byte, error) {
	return b.metricsMarshaler.MarshalMetrics(md)
}

func (b *baseMarshaler) Format() string {
	return "json"
}

// gzipMarshaler gzip compresses marshalled data
type gzipMarshaler struct {
	base *baseMarshaler
}

func (g *gzipMarshaler) MarshalTraces(td ptrace.Traces) ([]byte, error) {
	data, err := g.base.MarshalTraces(td)
	if err != nil {
		return nil, err
	}

	return g.compress(data)
}

func (g *gzipMarshaler) MarshalLogs(ld plog.Logs) ([]byte, error) {
	data, err := g.base.MarshalLogs(ld)
	if err != nil {
		return nil, err
	}

	return g.compress(data)
}

func (g *gzipMarshaler) MarshalMetrics(md pmetric.Metrics) ([]byte, error) {
	data, err := g.base.MarshalMetrics(md)
	if err != nil {
		return nil, err
	}

	return g.compress(data)
}

func (g *gzipMarshaler) Format() string {
	return "json.gz"
}

func (g *gzipMarshaler) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
