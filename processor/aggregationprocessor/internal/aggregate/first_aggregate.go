package aggregate

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type firstAggregate struct {
	firstValInt    int64
	firstValDouble float64
	firstTimestamp time.Time
	isInt          bool
}

func newFirstAggregate(initialVal pmetric.NumberDataPoint) (Aggregate, error) {
	switch initialVal.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		return &firstAggregate{
			firstValInt:    initialVal.IntValue(),
			isInt:          true,
			firstTimestamp: initialVal.Timestamp().AsTime(),
		}, nil
	case pmetric.NumberDataPointValueTypeDouble:
		return &firstAggregate{
			firstValDouble: initialVal.DoubleValue(),
			isInt:          false,
			firstTimestamp: initialVal.Timestamp().AsTime(),
		}, nil
	}

	return nil, fmt.Errorf("cannot create first aggregation from empty datapoint")
}

func (m *firstAggregate) AddDatapoint(ndp pmetric.NumberDataPoint) {
	if ndp.Timestamp() == 0 {
		// Ignore uninitialized timestamp
		return
	}

	ndpTimestamp := ndp.Timestamp().AsTime()
	if m.firstTimestamp.After(ndpTimestamp) {
		// ndp is before this metric, so we should use it first
		if m.isInt {
			m.firstValInt = getDatapointValueInt(ndp)
		} else {
			m.firstValDouble = getDatapointValueDouble(ndp)
		}
	}
}

func (m *firstAggregate) SetDatapointValue(dp pmetric.NumberDataPoint) {
	if m.isInt {
		dp.SetIntValue(m.firstValInt)
	} else {
		dp.SetDoubleValue(m.firstValDouble)
	}
}
