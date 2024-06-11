package bindplaneextension

import (
	"github.com/observiq/bindplane-agent/internal/measurements"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/otel/attribute"
)

const (
	reportMeasurementsCapability = "com.bindplane.measurements.v1"
	reportMeasurementsType       = "reportMeasurements"
)

func otlpMeasurements(tm *measurements.ThroughputMeasurements) pmetric.MetricSlice {
	s := pmetric.NewMetricSlice()

	attrs := pcommon.NewMap()
	sdkAttrs := tm.Attributes()
	attrIter := sdkAttrs.Iter()

	for attrIter.Next() {
		kv := attrIter.Attribute()
		switch kv.Value.Type() {
		case attribute.STRING:
			attrs.PutStr(string(kv.Key), kv.Value.AsString())
		default: // Do nothing for non-string attributes; Attributes for throughput metrics can only be strings for now.
		}
	}

	setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_log_data_size", tm.LogSize(), attrs)
	setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_metric_data_size", tm.MetricSize(), attrs)
	setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_trace_data_size", tm.TraceSize(), attrs)

	setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_log_count", tm.LogCount(), attrs)
	setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_metric_count", tm.DatapointCount(), attrs)
	setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_trace_count", tm.TraceSize(), attrs)

	return s
}

func setOTLPSum(m pmetric.Metric, name string, value int64, attrs pcommon.Map) {
	m.SetName(name)
	m.SetEmptySum()
	s := m.Sum()

	dp := s.DataPoints().AppendEmpty()
	dp.SetIntValue(value)
	attrs.CopyTo(dp.Attributes())
}
