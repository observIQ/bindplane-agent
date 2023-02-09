package aggregate

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type minAggregation struct {
	minDouble float64
	minInt    int64
	isInt     bool
}

func newMinAggregate(initialVal pmetric.NumberDataPoint) (Aggregate, error) {
	switch initialVal.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		return &minAggregation{
			minInt: initialVal.IntValue(),
			isInt:  true,
		}, nil
	case pmetric.NumberDataPointValueTypeDouble:
		return &minAggregation{
			minDouble: initialVal.DoubleValue(),
			isInt:     false,
		}, nil
	}

	return nil, fmt.Errorf("cannot create min aggregation from empty datapoint")
}

func (m *minAggregation) AddDatapoint(ndp pmetric.NumberDataPoint) {
	if m.isInt {
		i := getDatapointValueInt(ndp)
		if i < m.minInt {
			m.minInt = i
		}
	} else {
		f := getDatapointValueDouble(ndp)
		if f < m.minDouble {
			m.minDouble = f
		}
	}
}

func (m *minAggregation) SetDatapointValue(dp pmetric.NumberDataPoint) {
	if m.isInt {
		dp.SetIntValue(m.minInt)
	} else {
		dp.SetDoubleValue(m.minDouble)
	}
}
