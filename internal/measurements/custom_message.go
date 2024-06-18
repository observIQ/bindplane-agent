package measurements

import (
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/otel/attribute"
)

const (
	ReportMeasurementsV1Capability = "com.bindplane.measurements.v1"
	ReportMeasurementsType         = "reportMeasurements"
)

func OTLPThroughputMeasurements(tm *ThroughputMeasurements, includeCountMetrics bool) pmetric.MetricSlice {
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

	ts := pcommon.NewTimestampFromTime(time.Now())

	setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_log_data_size", tm.LogSize(), attrs, ts)
	setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_metric_data_size", tm.MetricSize(), attrs, ts)
	setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_trace_data_size", tm.TraceSize(), attrs, ts)

	if includeCountMetrics {
		setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_log_count", tm.LogCount(), attrs, ts)
		setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_metric_count", tm.DatapointCount(), attrs, ts)
		setOTLPSum(s.AppendEmpty(), "otelcol_processor_throughputmeasurement_trace_count", tm.TraceSize(), attrs, ts)
	}

	return s
}

func setOTLPSum(m pmetric.Metric, name string, value int64, attrs pcommon.Map, now pcommon.Timestamp) {
	if value == 0 {
		// Ignore value if it's 0
		return
	}

	m.SetName(name)
	m.SetEmptySum()
	s := m.Sum()

	dp := s.DataPoints().AppendEmpty()
	dp.SetIntValue(value)
	attrs.CopyTo(dp.Attributes())
	dp.SetTimestamp(now)
}
