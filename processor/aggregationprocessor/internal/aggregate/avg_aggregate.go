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

package aggregate

import (
	"errors"

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

	return nil, errors.New("cannot create avg aggregation from empty datapoint")
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
