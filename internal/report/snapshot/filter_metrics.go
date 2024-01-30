package snapshot

import (
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// filterMetrics filters the metrics by the given query and timestamp.
// The returned payload cannot be assumed to be a copy, so it should not be modified.
func filterMetrics(m pmetric.Metrics, query *string, minTimestamp *time.Time) pmetric.Metrics {
	// No filters specified, filtered metrics are trivially the same as input metrics
	if query == nil && minTimestamp == nil {
		return m
	}

	filteredMetrics := pmetric.NewMetrics()
	resourceMetrics := m.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		filteredResourceLogs := filterResourceMetrics(resourceMetrics.At(i), query, minTimestamp)

		// Don't append empty resource metrics
		if filteredResourceLogs.ScopeMetrics().Len() != 0 {
			filteredResourceLogs.MoveTo(filteredMetrics.ResourceMetrics().AppendEmpty())
		}
	}

	return filteredMetrics
}

func filterResourceMetrics(rm pmetric.ResourceMetrics, query *string, minTimestamp *time.Time) pmetric.ResourceMetrics {
	filteredResourceMetrics := pmetric.NewResourceMetrics()

	// Copy old resource to filtered resource
	resource := rm.Resource()
	resource.CopyTo(filteredResourceMetrics.Resource())

	// Apply query to resource
	queryMatchesResource := true // default to true if no query specified
	if query != nil {
		queryMatchesResource = queryMatchesMap(resource.Attributes(), *query)
	}

	scopeMetrics := rm.ScopeMetrics()
	for i := 0; i < scopeMetrics.Len(); i++ {
		filteredScopeMetrics := filterScopeMetrics(rm.ScopeMetrics().At(i), queryMatchesResource, query, minTimestamp)

		// Don't append empty scope metrics
		if filteredScopeMetrics.Metrics().Len() != 0 {
			filteredScopeMetrics.MoveTo(filteredResourceMetrics.ScopeMetrics().AppendEmpty())
		}
	}

	return filteredResourceMetrics
}

func filterScopeMetrics(sm pmetric.ScopeMetrics, queryMatchesResource bool, query *string, minTimestamp *time.Time) pmetric.ScopeMetrics {
	filteredScopeMetrics := pmetric.NewScopeMetrics()
	metrics := sm.Metrics()
	for i := 0; i < metrics.Len(); i++ {
		m := metrics.At(i)
		filteredMetric := filterMetric(m, queryMatchesResource, query, minTimestamp)

		if !metricIsEmpty(filteredMetric) {
			filteredMetric.MoveTo(filteredScopeMetrics.Metrics().AppendEmpty())
		}
	}

	return filteredScopeMetrics
}

func filterMetric(m pmetric.Metric, queryMatchesResource bool, query *string, minTimestamp *time.Time) pmetric.Metric {
	filteredMetric := pmetric.NewMetric()
	// Copy metric to filtered metric
	filteredMetric.SetName(m.Name())
	filteredMetric.SetDescription(m.Description())
	filteredMetric.SetUnit(m.Unit())

	// Apply query to metric
	queryMatchesMetric := true // default to true if no query specified
	// Skip if we already know the query matches the resource
	if !queryMatchesResource && query != nil {
		queryMatchesMetric = metricMatchesQuery(m, *query)
	}

	switch m.Type() {
	case pmetric.MetricTypeGauge:
		filteredGauge := filterGauge(m.Gauge(), queryMatchesResource, queryMatchesMetric, query, minTimestamp)
		filteredGauge.MoveTo(filteredMetric.SetEmptyGauge())
	case pmetric.MetricTypeSum:
		filteredSum := filterSum(m.Sum(), queryMatchesResource, queryMatchesMetric, query, minTimestamp)
		filteredSum.MoveTo(filteredMetric.SetEmptySum())
	case pmetric.MetricTypeHistogram:
		filteredHistogram := filterHistogram(m.Histogram(), queryMatchesResource, queryMatchesMetric, query, minTimestamp)
		filteredHistogram.MoveTo(filteredMetric.SetEmptyHistogram())
	case pmetric.MetricTypeExponentialHistogram:
		filteredExponentialHistogram := filterExponentialHistogram(m.ExponentialHistogram(), queryMatchesResource, queryMatchesMetric, query, minTimestamp)
		filteredExponentialHistogram.MoveTo(filteredMetric.SetEmptyExponentialHistogram())
	case pmetric.MetricTypeSummary:
		filteredSummary := filterSummary(m.Summary(), queryMatchesResource, queryMatchesMetric, query, minTimestamp)
		filteredSummary.MoveTo(filteredMetric.SetEmptySummary())
	case pmetric.MetricTypeEmpty:
		// Ignore empty
	}

	return filteredMetric
}

func filterGauge(g pmetric.Gauge, queryMatchesResource, queryMatchesName bool, query *string, minTimestamp *time.Time) pmetric.Gauge {
	filteredGauge := pmetric.NewGauge()

	dps := g.DataPoints()
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		if datapointMatches(dp, queryMatchesResource, queryMatchesName, query, minTimestamp) {
			dp.CopyTo(filteredGauge.DataPoints().AppendEmpty())
		}
	}

	return filteredGauge
}

func filterSum(s pmetric.Sum, queryMatchesResource, queryMatchesName bool, query *string, minTimestamp *time.Time) pmetric.Sum {
	filteredSum := pmetric.NewSum()

	dps := s.DataPoints()
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		if datapointMatches(dp, queryMatchesResource, queryMatchesName, query, minTimestamp) {
			dp.CopyTo(filteredSum.DataPoints().AppendEmpty())
		}
	}

	return filteredSum
}

func filterHistogram(h pmetric.Histogram, queryMatchesResource, queryMatchesName bool, query *string, minTimestamp *time.Time) pmetric.Histogram {
	filteredHistogram := pmetric.NewHistogram()

	dps := h.DataPoints()
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		if datapointMatches(dp, queryMatchesResource, queryMatchesName, query, minTimestamp) {
			dp.CopyTo(filteredHistogram.DataPoints().AppendEmpty())
		}
	}

	return filteredHistogram
}

func filterExponentialHistogram(eh pmetric.ExponentialHistogram, queryMatchesResource, queryMatchesName bool, query *string, minTimestamp *time.Time) pmetric.ExponentialHistogram {
	filteredExponentialHistogram := pmetric.NewExponentialHistogram()

	dps := eh.DataPoints()
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		if datapointMatches(dp, queryMatchesResource, queryMatchesName, query, minTimestamp) {
			dp.CopyTo(filteredExponentialHistogram.DataPoints().AppendEmpty())
		}
	}

	return filteredExponentialHistogram
}

func filterSummary(s pmetric.Summary, queryMatchesResource, queryMatchesName bool, query *string, minTimestamp *time.Time) pmetric.Summary {
	filteredSummary := pmetric.NewSummary()

	dps := s.DataPoints()
	for i := 0; i < dps.Len(); i++ {
		dp := dps.At(i)
		if datapointMatches(dp, queryMatchesResource, queryMatchesName, query, minTimestamp) {
			dp.CopyTo(filteredSummary.DataPoints().AppendEmpty())
		}
	}

	return filteredSummary
}

func metricMatchesQuery(m pmetric.Metric, query string) bool {
	// Match query against metric name
	return strings.Contains(m.Name(), query)
}

// datapoint is an interface that every concrete datapoint type implements
type datapoint interface {
	Attributes() pcommon.Map
	Timestamp() pcommon.Timestamp
}

func datapointMatches(dp datapoint, queryMatchesResource, queryMatchesName bool, query *string, minTimestamp *time.Time) bool {
	queryAlreadyMatched := queryMatchesResource || queryMatchesName

	queryMatchesDatapoint := true
	if !queryAlreadyMatched && query != nil {
		queryMatchesDatapoint = datapointMatchesQuery(dp, *query)
	}

	matchesTimestamp := true
	if minTimestamp != nil {
		matchesTimestamp = datapointMatchesTimestamp(dp, *minTimestamp)
	}

	matchesQuery := queryMatchesResource || queryMatchesName || queryMatchesDatapoint

	return matchesQuery && matchesTimestamp
}

func datapointMatchesTimestamp(dp datapoint, minTimestamp time.Time) bool {
	return dp.Timestamp() >= pcommon.NewTimestampFromTime(minTimestamp)
}

func datapointMatchesQuery(dp datapoint, query string) bool {
	return queryMatchesMap(dp.Attributes(), query)
}

func metricIsEmpty(m pmetric.Metric) bool {
	switch m.Type() {
	case pmetric.MetricTypeGauge:
		return m.Gauge().DataPoints().Len() == 0
	case pmetric.MetricTypeSum:
		return m.Sum().DataPoints().Len() == 0
	case pmetric.MetricTypeHistogram:
		return m.Histogram().DataPoints().Len() == 0
	case pmetric.MetricTypeExponentialHistogram:
		return m.ExponentialHistogram().DataPoints().Len() == 0
	case pmetric.MetricTypeSummary:
		return m.Summary().DataPoints().Len() == 0
	case pmetric.MetricTypeEmpty:
		return true
	}
	return false
}
