package aggregate

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

// Aggregate is an interface represents an aggregate of datapoints
type Aggregate interface {
	AddDatapoint(f pmetric.NumberDataPoint)
	SetDatapointValue(pmetric.NumberDataPoint)
}

// AggregationType represents a type of aggregate
type AggregationType string

// Types of aggregates
const (
	AggregationTypeMin   AggregationType = "min"
	AggregationTypeMax   AggregationType = "max"
	AggregationTypeFirst AggregationType = "first"
	AggregationTypeLast  AggregationType = "last"
	AggregationTypeAvg   AggregationType = "avg"
)

type aggregateConstructor func(pmetric.NumberDataPoint) (Aggregate, error)

var aggregateConstructors = map[AggregationType]aggregateConstructor{
	AggregationTypeMin:   newMinAggregate,
	AggregationTypeMax:   newMaxAggregate,
	AggregationTypeFirst: newFirstAggregate,
	AggregationTypeLast:  newLastAggregate,
	AggregationTypeAvg:   newAvgAggregate,
}

// New creates a new aggregate of the given type, using the initial datapoint
func (a AggregationType) New(initialVal pmetric.NumberDataPoint) (Aggregate, error) {
	constructor, ok := aggregateConstructors[a]
	if !ok {
		return nil, fmt.Errorf("invalid aggregation type: %s", a)
	}
	return constructor(initialVal)
}

// Valid returns true if this Type is a valid aggregate type, false otherwise
func (a AggregationType) Valid() bool {
	_, ok := aggregateConstructors[a]
	return ok
}
