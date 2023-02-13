package aggregationprocessor

import "go.opentelemetry.io/collector/pdata/pmetric"

// removeEmptyMetrics removes empty gauge or sum metrics that have no datapoints remaining
func removeEmptyMetrics(ms pmetric.MetricSlice) {
	ms.RemoveIf(func(m pmetric.Metric) bool {
		switch m.Type() {
		case pmetric.MetricTypeGauge:
			return m.Gauge().DataPoints().Len() == 0
		case pmetric.MetricTypeSum:
			return m.Sum().DataPoints().Len() == 0
		}
		return false
	})
}

// removeEmptyScopeMetrics removes any empty ScopeMetrics
func removeEmptyScopeMetrics(sms pmetric.ScopeMetricsSlice) {
	sms.RemoveIf(func(sm pmetric.ScopeMetrics) bool {
		return sm.Metrics().Len() == 0
	})
}

// removeEmptyResourceMetrics removes any empty ResourceMetrics
func removeEmptyResourceMetrics(rms pmetric.ResourceMetricsSlice) {
	rms.RemoveIf(func(rm pmetric.ResourceMetrics) bool {
		return rm.ScopeMetrics().Len() == 0
	})
}

// isMonotonic returns true if the metric is a monotonic sum, false otherwise.
func isMonotonic(m pmetric.Metric) bool {
	if m.Type() == pmetric.MetricTypeSum {
		return m.Sum().IsMonotonic()
	}
	// Monotonicity is only an attribute of the Sum type.
	return false
}

// datapointsFromMetric gets the underlying datapoint slice from gauge or sum metrics.
func datapointsFromMetric(m pmetric.Metric) pmetric.NumberDataPointSlice {
	switch m.Type() {
	case pmetric.MetricTypeGauge:
		return m.Gauge().DataPoints()
	case pmetric.MetricTypeSum:
		return m.Sum().DataPoints()
	}

	// getting datapoints only supported for gauge and sum types
	return pmetric.NewNumberDataPointSlice()
}
