package aggregate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestAggregateTypesValid(t *testing.T) {
	types := []AggregationType{
		AggregationTypeMin,
		AggregationTypeAvg,
		AggregationTypeMax,
		AggregationTypeFirst,
		AggregationTypeLast,
	}

	for _, aggType := range types {
		t.Run(string(aggType), func(t *testing.T) {
			require.True(t, aggType.Valid(), "Invalid aggregation type %s", aggType)
		})
	}
}

func TestAggregateTypesNew(t *testing.T) {
	types := []AggregationType{
		AggregationTypeMin,
		AggregationTypeAvg,
		AggregationTypeMax,
		AggregationTypeFirst,
		AggregationTypeLast,
	}

	for _, aggType := range types {
		t.Run(string(aggType), func(t *testing.T) {
			dp := pmetric.NewNumberDataPoint()
			dp.SetDoubleValue(2.0)
			a, err := aggType.New(dp)
			require.NoError(t, err)
			require.NotNil(t, a)
		})
	}
}

func TestAggregateTypeNewInvalidType(t *testing.T) {
	dp := pmetric.NewNumberDataPoint()
	dp.SetDoubleValue(2.0)
	_, err := AggregationType("invalid").New(dp)
	require.ErrorContains(t, err, "invalid aggregation type:")
}

func TestAggregateTypesNewEmptyDatapoint(t *testing.T) {
	types := []AggregationType{
		AggregationTypeMin,
		AggregationTypeAvg,
		AggregationTypeMax,
		AggregationTypeFirst,
		AggregationTypeLast,
	}

	for _, aggType := range types {
		t.Run(string(aggType), func(t *testing.T) {
			dp := pmetric.NewNumberDataPoint()
			_, err := aggType.New(dp)
			require.ErrorContains(t, err, fmt.Sprintf("cannot create %s aggregation from empty datapoint", aggType))
		})
	}
}

func TestAggregatesFloat(t *testing.T) {
	testCases := []struct {
		name       string
		aggType    AggregationType
		values     []float64
		timestamps []int64
		finalValue float64
	}{
		{
			name:       "min",
			aggType:    AggregationTypeMin,
			values:     []float64{45, 1, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 1,
		},
		{
			name:       "max",
			aggType:    AggregationTypeMax,
			values:     []float64{45, 1, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 99,
		},
		{
			name:       "avg",
			aggType:    AggregationTypeAvg,
			values:     []float64{45, 0, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 36.75,
		},
		{
			name:       "first (unset timestamp)",
			aggType:    AggregationTypeFirst,
			values:     []float64{45, 1, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 45,
		},
		{
			name:       "last (unset timestamp)",
			aggType:    AggregationTypeLast,
			values:     []float64{45, 1, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 3,
		},
		{
			name:       "first (set timestamp)",
			aggType:    AggregationTypeFirst,
			values:     []float64{45, 1, 99, 3},
			timestamps: []int64{10, 2, 3, 11},
			finalValue: 99,
		},
		{
			name:       "last (set timestamp)",
			aggType:    AggregationTypeLast,
			values:     []float64{45, 1, 99, 3},
			timestamps: []int64{10, 3, 89, 11},
			finalValue: 99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialVal := pmetric.NewNumberDataPoint()
			initialVal.SetDoubleValue(tc.values[0])
			initialVal.SetTimestamp(pcommon.Timestamp(tc.timestamps[0]))
			agg, err := tc.aggType.New(initialVal)
			require.NoError(t, err)

			for i, v := range tc.values[1:] {
				dp := pmetric.NewNumberDataPoint()
				dp.SetTimestamp(pcommon.Timestamp(tc.timestamps[i+1]))
				dp.SetDoubleValue(v)
				agg.AddDatapoint(dp)
			}
			finalDp := pmetric.NewNumberDataPoint()
			agg.SetDatapointValue(finalDp)

			require.Equal(t, tc.finalValue, finalDp.DoubleValue())
		})
	}
}

func TestAggregatesInt(t *testing.T) {
	testCases := []struct {
		name       string
		aggType    AggregationType
		values     []int64
		timestamps []int64
		finalValue int64
	}{
		{
			name:       "min",
			aggType:    AggregationTypeMin,
			values:     []int64{45, 1, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 1,
		},
		{
			name:       "max",
			aggType:    AggregationTypeMax,
			values:     []int64{45, 1, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 99,
		},
		{
			name:       "avg",
			aggType:    AggregationTypeAvg,
			values:     []int64{45, 0, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 36,
		},
		{
			name:       "first (unset timestamp)",
			aggType:    AggregationTypeFirst,
			values:     []int64{45, 1, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 45,
		},
		{
			name:       "last (unset timestamp)",
			aggType:    AggregationTypeLast,
			values:     []int64{45, 1, 99, 3},
			timestamps: []int64{0, 0, 0, 0},
			finalValue: 3,
		},
		{
			name:       "first (set timestamp)",
			aggType:    AggregationTypeFirst,
			values:     []int64{45, 1, 99, 3},
			timestamps: []int64{10, 2, 3, 11},
			finalValue: 99,
		},
		{
			name:       "last (set timestamp)",
			aggType:    AggregationTypeLast,
			values:     []int64{45, 1, 99, 3},
			timestamps: []int64{10, 3, 89, 11},
			finalValue: 99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			initialVal := pmetric.NewNumberDataPoint()
			initialVal.SetIntValue(tc.values[0])
			initialVal.SetTimestamp(pcommon.Timestamp(tc.timestamps[0]))
			agg, err := tc.aggType.New(initialVal)
			require.NoError(t, err)

			for i, v := range tc.values[1:] {
				dp := pmetric.NewNumberDataPoint()
				dp.SetTimestamp(pcommon.Timestamp(tc.timestamps[i+1]))
				dp.SetIntValue(v)
				agg.AddDatapoint(dp)
			}
			finalDp := pmetric.NewNumberDataPoint()
			agg.SetDatapointValue(finalDp)

			require.Equal(t, tc.finalValue, finalDp.IntValue())
		})
	}
}
