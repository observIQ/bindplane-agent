package logsummaryprocessor

import (
	"encoding/json"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// summary is a summary of observed logs.
type summary struct {
	resources map[string]*resourceSummary
}

// newSummary returns a new summary.
func newSummary() *summary {
	return &summary{
		resources: make(map[string]*resourceSummary),
	}
}

// toMetrics returns the summary as a colleciton of metrics.
func (s *summary) toMetrics() pmetric.Metrics {
	metrics := pmetric.NewMetrics()
	for _, resource := range s.resources {
		resourceMetrics := metrics.ResourceMetrics().AppendEmpty()
		resource.attributes.CopyTo(resourceMetrics.Resource().Attributes())

		scopeMetrics := resourceMetrics.ScopeMetrics().AppendEmpty()
		scopeMetrics.Scope().SetName("logsummary")

		for _, logs := range resource.logs {
			logMetrics := scopeMetrics.Metrics().AppendEmpty()
			logMetrics.SetName("log.count")
			logMetrics.SetUnit("{logs}")
			gauge := logMetrics.Gauge().DataPoints().AppendEmpty()
			gauge.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
			gauge.SetIntValue(int64(logs.count))
		}
	}
	return metrics
}

// reset resets the summary.
func (s *summary) reset() {
	for k := range s.resources {
		delete(s.resources, k)
	}
}

// update will update the summary.
func (s *summary) update(pl plog.Logs) error {
	resLogSlice := pl.ResourceLogs()
	for i := 0; i < resLogSlice.Len(); i++ {
		resLogs := resLogSlice.At(i)

		attrs := resLogs.Resource().Attributes()
		attrsJSON, err := json.Marshal(attrs.AsRaw())
		if err != nil {
			return err
		}
		resourceKey := string(attrsJSON)

		if _, ok := s.resources[resourceKey]; !ok {
			s.resources[resourceKey] = newResourceSummary()
			attrs.CopyTo(s.resources[resourceKey].attributes)
		}

		if err := s.resources[resourceKey].update(resLogs); err != nil {
			return err
		}
	}

	return nil
}

// resourceSummary is a summary of resource logs observed.
type resourceSummary struct {
	attributes pcommon.Map
	logs       map[string]*logSummary
}

// newResourceSummary returns a new resource summary.
func newResourceSummary() *resourceSummary {
	return &resourceSummary{
		attributes: pcommon.NewMap(),
		logs:       make(map[string]*logSummary),
	}
}

// update will update the resource summary.
func (r *resourceSummary) update(resLogs plog.ResourceLogs) error {
	for j := 0; j < resLogs.ScopeLogs().Len(); j++ {
		scopeLog := resLogs.ScopeLogs().At(j)
		logRecords := scopeLog.LogRecords()
		for k := 0; k < logRecords.Len(); k++ {
			logRecord := logRecords.At(k)
			attrs := pcommon.NewMap()
			logRecord.Attributes().CopyTo(attrs)
			attrs.PutString("severity", logRecord.SeverityText())

			attrsJSON, err := json.Marshal(attrs.AsRaw())
			if err != nil {
				return err
			}
			recordKey := string(attrsJSON)

			if _, ok := r.logs[recordKey]; !ok {
				r.logs[recordKey] = newLogSummary()
				attrs.CopyTo(r.logs[recordKey].attributes)
			}

			r.logs[recordKey].count++
			return nil
		}
	}

	return nil
}

// logSummary is a summary of log records observed.
type logSummary struct {
	attributes pcommon.Map
	count      int
}

// newLogSummary returns a new log summary.
func newLogSummary() *logSummary {
	return &logSummary{
		attributes: pcommon.NewMap(),
	}
}
