package aggregate

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type avgAggregate struct {
	totalInt    int64
	totalDouble float64
	isInt       bool
	count       int64
}

func newAvgAggregate(initialVal pmetric.NumberDataPoint) (Aggregate, error) {
	switch initialVal.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		return &avgAggregate{
			totalInt: initialVal.IntValue(),
			isInt:    true,
			count:    1,
		}, nil
	case pmetric.NumberDataPointValueTypeDouble:
		return &avgAggregate{
			totalDouble: initialVal.DoubleValue(),
			isInt:       false,
			count:       1,
		}, nil
	}

	return nil, fmt.Errorf("cannot create avg aggregation from empty datapoint")
}

func (m *avgAggregate) AddDatapoint(ndp pmetric.NumberDataPoint) {
	if m.isInt {
		i := getDatapointValueInt(ndp)
		m.totalInt += i
	} else {
		f := getDatapointValueDouble(ndp)
		m.totalDouble += f
	}

	m.count++
}

func (m *avgAggregate) SetDatapointValue(dp pmetric.NumberDataPoint) {
	if m.isInt {
		dp.SetIntValue(m.totalInt / m.count)
	} else {
		dp.SetDoubleValue(m.totalDouble / float64(m.count))
	}
}
