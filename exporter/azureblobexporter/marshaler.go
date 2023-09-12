package azureblobexporter // import "github.com/observiq/bindplane-agent/exporter/azureblobexporter"

import (
	"bytes"
	"compress/gzip"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// marshaler marshals into data bytes based on configuration
type marshaler interface {
	MarshalTraces(td ptrace.Traces) ([]byte, error)
	MarshalLogs(ld plog.Logs) ([]byte, error)
	MarshalMetrics(md pmetric.Metrics) ([]byte, error)
	format() string
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

func (b *baseMarshaler) format() string {
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

func (g *gzipMarshaler) format() string {
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
