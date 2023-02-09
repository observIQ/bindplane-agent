package aggregate

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type lastAggregation struct {
	lastValInt    int64
	lastValDouble float64
	lastTimestamp time.Time
	isInt         bool
}

func newLastAggregate(initialVal pmetric.NumberDataPoint) (Aggregate, error) {
	switch initialVal.ValueType() {
	case pmetric.NumberDataPointValueTypeInt:
		return &lastAggregation{
			lastValInt:    initialVal.IntValue(),
			isInt:         true,
			lastTimestamp: initialVal.Timestamp().AsTime(),
		}, nil
	case pmetric.NumberDataPointValueTypeDouble:
		return &lastAggregation{
			lastValDouble: initialVal.DoubleValue(),
			isInt:         false,
			lastTimestamp: initialVal.Timestamp().AsTime(),
		}, nil
	}

	return nil, fmt.Errorf("cannot create last aggregation from empty datapoint")
}

func (m *lastAggregation) AddDatapoint(ndp pmetric.NumberDataPoint) {
	ndpTimestamp := ndp.Timestamp().AsTime()
	// Note for the Equal here: If two datapoints have the same timestamp, we consider the last
	// datapoint we receive from the pipeline to be the "last" datapoint.
	if ndpTimestamp.After(m.lastTimestamp) || ndpTimestamp.Equal(m.lastTimestamp) {
		// ndp is after this metric, so we should use it as the last value
		if m.isInt {
			m.lastValInt = getDatapointValueInt(ndp)
		} else {
			m.lastValDouble = getDatapointValueDouble(ndp)
		}

		m.lastTimestamp = ndpTimestamp
	}
}

func (m *lastAggregation) SetDatapointValue(dp pmetric.NumberDataPoint) {
	if m.isInt {
		dp.SetIntValue(m.lastValInt)
	} else {
		dp.SetDoubleValue(m.lastValDouble)
	}
}
