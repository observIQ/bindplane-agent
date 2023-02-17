// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package stats implements structs that are used to calculate statistics from datapoints.
package stats

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

// Statistic is an interface represents an aggregate of datapoints
type Statistic interface {
	AddDatapoint(f pmetric.NumberDataPoint)
	SetDatapointValue(pmetric.NumberDataPoint)
}

// StatType represents a type of statistic to calculate
type StatType string

// Types of aggregates
const (
	MinType   StatType = "min"
	MaxType   StatType = "max"
	FirstType StatType = "first"
	LastType  StatType = "last"
	AvgType   StatType = "avg"
)

type statConstructor func(pmetric.NumberDataPoint) (Statistic, error)

var statConstructors = map[StatType]statConstructor{
	MinType:   newMinStatistic,
	MaxType:   newMaxStatistic,
	FirstType: newFirstStatistic,
	LastType:  newLastStatistic,
	AvgType:   newAvgStatistic,
}

// New creates a new aggregate of the given type, using the initial datapoint
func (a StatType) New(initialVal pmetric.NumberDataPoint) (Statistic, error) {
	constructor, ok := statConstructors[a]
	if !ok {
		return nil, fmt.Errorf("invalid aggregation type: %s", a)
	}
	return constructor(initialVal)
}

// Valid returns true if this Type is a valid aggregate type, false otherwise
func (a StatType) Valid() bool {
	_, ok := statConstructors[a]
	return ok
}
