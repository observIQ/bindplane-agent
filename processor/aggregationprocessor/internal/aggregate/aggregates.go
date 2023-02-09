package aggregate

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type Aggregate interface {
	AddDatapoint(f pmetric.NumberDataPoint)
	SetDatapointValue(pmetric.NumberDataPoint)
}

type AggregateType string

const (
	AggregateTypeMin   AggregateType = "min"
	AggregateTypeMax   AggregateType = "max"
	AggregateTypeFirst AggregateType = "first"
	AggregateTypeLast  AggregateType = "last"
	AggregateTypeAvg   AggregateType = "avg"
)

type aggregateConstructor func(pmetric.NumberDataPoint) (Aggregate, error)

var aggregateConstructors = map[AggregateType]aggregateConstructor{
	AggregateTypeMin:   newMinAggregate,
	AggregateTypeMax:   newMaxAggregate,
	AggregateTypeFirst: newFirstAggregate,
	AggregateTypeLast:  newLastAggregate,
	AggregateTypeAvg:   newAvgAggregate,
}

func (a AggregateType) New(initialVal pmetric.NumberDataPoint) (Aggregate, error) {
	constructor, ok := aggregateConstructors[a]
	if !ok {
		return nil, fmt.Errorf("invalid aggregation type: %s", a)
	}
	return constructor(initialVal)
}

func (a AggregateType) Valid() bool {
	_, ok := aggregateConstructors[a]
	return ok
}
