package azureblobexporter

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
	base := &baseMarshaller{
		logsMarshaler:    &plog.JSONMarshaler{},
		tracesMarshaler:  &ptrace.JSONMarshaler{},
		metricsMarshaler: &pmetric.JSONMarshaler{},
	}

	switch compression {
	case gzipCompression:
		return &gzipMarshaller{
			base: base,
		}
	default:
		return base
	}
}

// baseMarshaller is the base marshaller that marshals otlp structures into JSON
type baseMarshaller struct {
	logsMarshaler    plog.Marshaler
	tracesMarshaler  ptrace.Marshaler
	metricsMarshaler pmetric.Marshaler
}

func (b *baseMarshaller) MarshalTraces(td ptrace.Traces) ([]byte, error) {
	return b.tracesMarshaler.MarshalTraces(td)
}

func (b *baseMarshaller) MarshalLogs(ld plog.Logs) ([]byte, error) {
	return b.logsMarshaler.MarshalLogs(ld)
}

func (b *baseMarshaller) MarshalMetrics(md pmetric.Metrics) ([]byte, error) {
	return b.metricsMarshaler.MarshalMetrics(md)
}

func (b *baseMarshaller) format() string {
	return "json"
}

// gzipMarshaller gzip compresses marshalled data
type gzipMarshaller struct {
	base *baseMarshaller
}

func (g *gzipMarshaller) MarshalTraces(td ptrace.Traces) ([]byte, error) {
	data, err := g.base.MarshalTraces(td)
	if err != nil {
		return nil, err
	}

	return g.compress(data)
}

func (g *gzipMarshaller) MarshalLogs(ld plog.Logs) ([]byte, error) {
	data, err := g.base.MarshalLogs(ld)
	if err != nil {
		return nil, err
	}

	return g.compress(data)
}

func (g *gzipMarshaller) MarshalMetrics(md pmetric.Metrics) ([]byte, error) {
	data, err := g.base.MarshalMetrics(md)
	if err != nil {
		return nil, err
	}

	return g.compress(data)
}

func (g *gzipMarshaller) format() string {
	return "json.gz"
}

func (g *gzipMarshaller) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
