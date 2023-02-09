package aggregate

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type maxAggregation struct {
	maxDouble float64
	maxInt    int64
	isInt     bool
}

func newMaxAggregate(initialVal pmetric.NumberDataPoint) (Aggregate, error) {
	switch initialVal.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		return &maxAggregation{
			maxInt: initialVal.IntValue(),
			isInt:  true,
		}, nil
	case pmetric.NumberDataPointValueTypeDouble:
		return &maxAggregation{
			maxDouble: initialVal.DoubleValue(),
			isInt:     false,
		}, nil
	}

	return nil, fmt.Errorf("cannot create max aggregation from empty datapoint")
}

func (m *maxAggregation) AddDatapoint(ndp pmetric.NumberDataPoint) {
	if m.isInt {
		i := getDatapointValueInt(ndp)
		if i > m.maxInt {
			m.maxInt = i
		}
	} else {
		f := getDatapointValueDouble(ndp)
		if f > m.maxDouble {
			m.maxDouble = f
		}
	}
}

func (m *maxAggregation) SetDatapointValue(dp pmetric.NumberDataPoint) {
	if m.isInt {
		dp.SetIntValue(m.maxInt)
	} else {
		dp.SetDoubleValue(m.maxDouble)
	}
}
