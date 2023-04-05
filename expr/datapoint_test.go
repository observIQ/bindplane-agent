package expr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestConvertToDatapointResourceGroup(t *testing.T) {
	now := time.Now().UTC()
	oneMinuteAgo := now.Add(-time.Minute)
	testResource1 := map[string]any{
		"resource": "attributes",
	}
	testResource2 := map[string]any{
		"resource": "attributes",
	}
	testAttrs := map[string]any{
		"attributes": "attributes",
	}

	metrics := pmetric.NewMetrics()
	resourceMetrics1 := metrics.ResourceMetrics().AppendEmpty()
	resourceMetrics1.Resource().Attributes().FromRaw(testResource1)

	resourceMetrics2 := metrics.ResourceMetrics().AppendEmpty()
	resourceMetrics2.Resource().Attributes().FromRaw(testResource2)

	metricSlice1 := resourceMetrics1.ScopeMetrics().AppendEmpty().Metrics()
	sumMetric(t, testAttrs, oneMinuteAgo, now, metricSlice1)
	gaugeMetric(t, testAttrs, oneMinuteAgo, now, metricSlice1)
	histogramMetric(t, testAttrs, oneMinuteAgo, now, metricSlice1)

	metricSlice2 := resourceMetrics2.ScopeMetrics().AppendEmpty().Metrics()
	exponentialHistogramMetric(t, testAttrs, oneMinuteAgo, now, metricSlice2)
	summaryMetric(t, testAttrs, oneMinuteAgo, now, metricSlice2)

	resourceGroups := ConvertToDatapointResourceGroup(metrics)

	require.Equal(t, []DatapointResourceGroup{
		{
			Resource: testResource1,
			Datapoints: []Datapoint{
				{
					MetricNameField:     "sum",
					ResourceField:       testResource1,
					AttributesField:     testAttrs,
					DatapointValueField: float64(300),
				},
				{
					MetricNameField:     "gauge",
					ResourceField:       testResource1,
					AttributesField:     testAttrs,
					DatapointValueField: int64(45),
				},
				{
					MetricNameField: "histogram",
					ResourceField:   testResource1,
					AttributesField: testAttrs,
				},
			},
		},
		{
			Resource: testResource2,
			Datapoints: []Datapoint{
				{
					MetricNameField: "exponential_histogram",
					ResourceField:   testResource1,
					AttributesField: testAttrs,
				},
				{
					MetricNameField: "summary",
					ResourceField:   testResource1,
					AttributesField: testAttrs,
				},
			},
		},
	}, resourceGroups)
}

func sumMetric(t *testing.T, attrs map[string]any, start, now time.Time, s pmetric.MetricSlice) {
	t.Helper()

	metric := s.AppendEmpty()
	metric.SetName("sum")
	dps := metric.SetEmptySum().DataPoints()
	dp := dps.AppendEmpty()

	require.NoError(t, dp.Attributes().FromRaw(attrs))
	dp.SetStartTimestamp(pcommon.NewTimestampFromTime(start))
	dp.SetTimestamp(pcommon.NewTimestampFromTime(now))
	dp.SetDoubleValue(300)
}

func gaugeMetric(t *testing.T, attrs map[string]any, start, now time.Time, s pmetric.MetricSlice) {
	t.Helper()

	metric := s.AppendEmpty()
	metric.SetName("gauge")
	dps := metric.SetEmptyGauge().DataPoints()
	dp := dps.AppendEmpty()

	require.NoError(t, dp.Attributes().FromRaw(attrs))
	dp.SetStartTimestamp(pcommon.NewTimestampFromTime(start))
	dp.SetTimestamp(pcommon.NewTimestampFromTime(now))
	dp.SetIntValue(45)
}

func histogramMetric(t *testing.T, attrs map[string]any, start, now time.Time, s pmetric.MetricSlice) {
	t.Helper()

	metric := s.AppendEmpty()
	metric.SetName("histogram")
	dps := metric.SetEmptyHistogram().DataPoints()
	dp := dps.AppendEmpty()

	require.NoError(t, dp.Attributes().FromRaw(attrs))
	dp.SetStartTimestamp(pcommon.NewTimestampFromTime(start))
	dp.SetTimestamp(pcommon.NewTimestampFromTime(now))
}

func exponentialHistogramMetric(t *testing.T, attrs map[string]any, start, now time.Time, s pmetric.MetricSlice) {
	t.Helper()

	metric := s.AppendEmpty()
	metric.SetName("exponential_histogram")
	dps := metric.SetEmptyExponentialHistogram().DataPoints()
	dp := dps.AppendEmpty()

	require.NoError(t, dp.Attributes().FromRaw(attrs))
	dp.SetStartTimestamp(pcommon.NewTimestampFromTime(start))
	dp.SetTimestamp(pcommon.NewTimestampFromTime(now))
}

func summaryMetric(t *testing.T, attrs map[string]any, start, now time.Time, s pmetric.MetricSlice) {
	t.Helper()

	metric := s.AppendEmpty()
	metric.SetName("summary")
	dps := metric.SetEmptySummary().DataPoints()
	dp := dps.AppendEmpty()

	require.NoError(t, dp.Attributes().FromRaw(attrs))
	dp.SetStartTimestamp(pcommon.NewTimestampFromTime(start))
	dp.SetTimestamp(pcommon.NewTimestampFromTime(now))
}

func emptyMetric(t *testing.T, s pmetric.MetricSlice) {
	t.Helper()

	metric := s.AppendEmpty()
	metric.SetName("empty_metric")
}
