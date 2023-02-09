package aggregationprocessor

import (
	"github.com/observiq/observiq-otel-collector/processor/aggregationprocessor/internal/aggregate"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type resourceMetadata struct {
	resource pcommon.Map
	// metric name -> metric aggregation
	metrics map[string]*metricMetadata
}

type metricMetadata struct {
	name       string
	desc       string
	unit       string
	metricType pmetric.MetricType
	// Only relevant to sum metrics
	monotonic bool
	// Map of attributes hash to datapointAggregation
	datapoints map[uint64]*datapointMetadata
}

type datapointMetadata struct {
	attributes pcommon.Map
	aggregates map[AggregateConfig]aggregate.Aggregate
}
