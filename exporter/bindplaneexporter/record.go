package bindplaneexporter

import (
	"time"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// MetricRecord is a metric record sent to bindplane
type MetricRecord struct {
	Name       string                 `json:"name"`
	Timestamp  time.Time              `json:"timestamp"`
	Value      interface{}            `json:"value"`
	Unit       string                 `json:"unit"`
	Type       string                 `json:"type"`
	Attributes map[string]interface{} `json:"attributes"`
	Resource   map[string]interface{} `json:"resource"`
}

// LogRecord is a log record sent to bindplane
type LogRecord struct {
	Timestamp  time.Time              `json:"timestamp"`
	Body       interface{}            `json:"body"`
	Severity   string                 `json:"severity"`
	Attributes map[string]interface{} `json:"attributes"`
	Resource   map[string]interface{} `json:"resource"`
}

// TraceRecord is a trace record sent to bindplane
type TraceRecord struct {
	Name         string                 `json:"name"`
	TraceID      string                 `json:"trace_id"`
	SpanID       string                 `json:"span_id"`
	ParentSpanID string                 `json:"parent_span_id"`
	Start        time.Time              `json:"start"`
	End          time.Time              `json:"end"`
	Attributes   map[string]interface{} `json:"attributes"`
	Resource     map[string]interface{} `json:"resource"`
}

// getRecordsFromMetrics gets metric records from pmetrics
func getRecordsFromMetrics(metrics pmetric.Metrics) []MetricRecord {
	records := []MetricRecord{}
	for i := 0; i < metrics.ResourceMetrics().Len(); i++ {
		resourceMetrics := metrics.ResourceMetrics().At(i)
		resourceAttributes := resourceMetrics.Resource().Attributes().AsRaw()
		for k := 0; k < resourceMetrics.ScopeMetrics().Len(); k++ {
			scopeMetrics := resourceMetrics.ScopeMetrics().At(k)
			for j := 0; j < scopeMetrics.Metrics().Len(); j++ {
				metric := scopeMetrics.Metrics().At(j)
				recordSlice := getRecordsFromMetric(metric, resourceAttributes)
				records = append(records, recordSlice...)
			}
		}
	}
	return records
}

// getRecords gets metric records from a pmetric
func getRecordsFromMetric(metric pmetric.Metric, resourceAttributes map[string]interface{}) []MetricRecord {
	switch metric.DataType() {
	case pmetric.MetricDataTypeSum:
		return getMetricRecordsFromSum(metric, resourceAttributes)
	case pmetric.MetricDataTypeGauge:
		return getMetricRecordsFromGauge(metric, resourceAttributes)
	case pmetric.MetricDataTypeSummary:
		return getMetricRecordsFromSummary(metric, resourceAttributes)
	}

	return nil
}

// getMetricRecordsFromSum converts a sum into metric records
func getMetricRecordsFromSum(metric pmetric.Metric, resourceAttributes map[string]interface{}) []MetricRecord {
	metricName := metric.Name()
	metricUnit := metric.Unit()
	metricType := metric.DataType().String()
	metricRecords := []MetricRecord{}
	sum := metric.Sum()
	points := sum.DataPoints()
	for i := 0; i < points.Len(); i++ {
		point := points.At(i)
		record := MetricRecord{
			Name:       metricName,
			Timestamp:  point.Timestamp().AsTime(),
			Value:      getDataPointValue(point),
			Unit:       metricUnit,
			Type:       metricType,
			Attributes: point.Attributes().AsRaw(),
			Resource:   resourceAttributes,
		}
		metricRecords = append(metricRecords, record)
	}
	return metricRecords
}

// getMetricRecordsFromGauge converts a gauge into metric records
func getMetricRecordsFromGauge(metric pmetric.Metric, resourceAttributes map[string]interface{}) []MetricRecord {
	metricName := metric.Name()
	metricUnit := metric.Unit()
	metricType := metric.DataType().String()
	metricRecords := []MetricRecord{}
	gauge := metric.Gauge()
	points := gauge.DataPoints()
	for i := 0; i < points.Len(); i++ {
		point := points.At(i)
		record := MetricRecord{
			Name:       metricName,
			Timestamp:  point.Timestamp().AsTime(),
			Value:      getDataPointValue(point),
			Unit:       metricUnit,
			Type:       metricType,
			Attributes: point.Attributes().AsRaw(),
			Resource:   resourceAttributes,
		}
		metricRecords = append(metricRecords, record)
	}
	return metricRecords
}

// getMetricRecordsFromSummary converts a summary into metric records
func getMetricRecordsFromSummary(metric pmetric.Metric, resourceAttributes map[string]interface{}) []MetricRecord {
	metricName := metric.Name()
	metricUnit := metric.Unit()
	metricType := metric.DataType().String()
	metricRecords := []MetricRecord{}
	summary := metric.Summary()
	points := summary.DataPoints()
	for i := 0; i < points.Len(); i++ {
		point := points.At(i)
		record := MetricRecord{
			Name:       metricName,
			Timestamp:  point.Timestamp().AsTime(),
			Value:      getSummaryPointValue(point),
			Unit:       metricUnit,
			Type:       metricType,
			Attributes: point.Attributes().AsRaw(),
			Resource:   resourceAttributes,
		}
		metricRecords = append(metricRecords, record)
	}
	return metricRecords
}

// getDataPointValue gets the value of a data point
func getDataPointValue(point pmetric.NumberDataPoint) interface{} {
	switch point.ValueType() {
	case pmetric.NumberDataPointValueTypeDouble:
		return point.DoubleVal()
	case pmetric.NumberDataPointValueTypeInt:
		return point.IntVal()
	default:
		return 0
	}
}

// getSummaryPointValue gets the value of a summary point
func getSummaryPointValue(point pmetric.SummaryDataPoint) map[float64]interface{} {
	value := make(map[float64]interface{})
	for i := 0; i < point.QuantileValues().Len(); i++ {
		q := point.QuantileValues().At(i)
		value[q.Quantile()] = q.Value()
	}
	return value
}

// getRecordsFromLogs gets log records from plogs
func getRecordsFromLogs(logs plog.Logs) []LogRecord {
	records := []LogRecord{}
	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		resourceLogs := logs.ResourceLogs().At(i)
		resourceAttributes := resourceLogs.Resource().Attributes().AsRaw()
		for k := 0; k < resourceLogs.ScopeLogs().Len(); k++ {
			scopeLogs := resourceLogs.ScopeLogs().At(k)
			for j := 0; j < scopeLogs.LogRecords().Len(); j++ {
				log := scopeLogs.LogRecords().At(j)
				record := LogRecord{
					Timestamp:  log.Timestamp().AsTime(),
					Body:       log.Body().AsString(),
					Severity:   log.SeverityText(),
					Attributes: log.Attributes().AsRaw(),
					Resource:   resourceAttributes,
				}
				records = append(records, record)
			}
		}
	}
	return records
}

// getRecordsFromTraces gets trace records from ptraces
func getRecordsFromTraces(traces ptrace.Traces) []TraceRecord {
	records := []TraceRecord{}
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		resourceSpans := traces.ResourceSpans().At(i)
		resourceAttributes := resourceSpans.Resource().Attributes().AsRaw()
		for k := 0; k < resourceSpans.ScopeSpans().Len(); k++ {
			scopeSpans := resourceSpans.ScopeSpans().At(k)
			for j := 0; j < scopeSpans.Spans().Len(); j++ {
				span := scopeSpans.Spans().At(j)
				span.TraceID().HexString()
				record := TraceRecord{
					Name:         span.Name(),
					TraceID:      span.TraceID().HexString(),
					SpanID:       span.SpanID().HexString(),
					ParentSpanID: span.ParentSpanID().HexString(),
					Start:        span.StartTimestamp().AsTime(),
					End:          span.EndTimestamp().AsTime(),
					Attributes:   span.Attributes().AsRaw(),
					Resource:     resourceAttributes,
				}
				records = append(records, record)
			}
		}
	}
	return records
}
