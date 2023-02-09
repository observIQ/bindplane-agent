package aggregate

import "go.opentelemetry.io/collector/pdata/pmetric"

func getDatapointValueDouble(ndp pmetric.NumberDataPoint) float64 {
	switch ndp.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		return float64(ndp.IntValue())
	case pmetric.NumberDataPointValueTypeDouble:
		return ndp.DoubleValue()
	}

	// Empty number datapoint, we'll just return 0 in this case.
	// It's up to the caller to handle this case correctly.
	return 0
}

func getDatapointValueInt(ndp pmetric.NumberDataPoint) int64 {
	switch ndp.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		return ndp.IntValue()
	case pmetric.NumberDataPointValueTypeDouble:
		return int64(ndp.DoubleValue())
	}

	// Empty number datapoint, we'll just return 0 in this case.
	// It's up to the caller to handle this case correctly.
	return 0
}
