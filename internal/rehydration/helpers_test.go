// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rehydration //import "github.com/observiq/bindplane-agent/internal/rehydration"

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
)

func TestParseEntityPath(t *testing.T) {
	expectedTimeMinute := time.Date(2023, time.January, 04, 12, 02, 0, 0, time.UTC)
	expectedTimeHour := time.Date(2023, time.January, 04, 12, 00, 0, 0, time.UTC)

	testcases := []struct {
		desc         string
		entityName   string
		expectedTime *time.Time
		expectedType component.DataType
		expectedErr  error
	}{
		{
			desc:         "Empty entityName",
			entityName:   "",
			expectedTime: nil,
			expectedType: component.Type{},
			expectedErr:  ErrInvalidEntityPath,
		},
		{
			desc:         "Malformed path",
			entityName:   "year=2023/day=04/hour=12/minute=02/entitymetrics_12345.json",
			expectedTime: nil,
			expectedType: component.Type{},
			expectedErr:  ErrInvalidEntityPath,
		},
		{
			desc:         "Malformed timestamp",
			entityName:   "year=2003/month=00/day=04/hour=12/minute=01/entitymetrics_12345.json",
			expectedTime: nil,
			expectedType: component.Type{},
			expectedErr:  errors.New("parse entity time"),
		},
		{
			desc:         "Prefix, minute, metrics",
			entityName:   "prefix/year=2023/month=01/day=04/hour=12/minute=02/entitymetrics_12345.json",
			expectedTime: &expectedTimeMinute,
			expectedType: component.DataTypeMetrics,
			expectedErr:  nil,
		},
		{
			desc:         "No Prefix, minute, metrics",
			entityName:   "year=2023/month=01/day=04/hour=12/minute=02/entitymetrics_12345.json",
			expectedTime: &expectedTimeMinute,
			expectedType: component.DataTypeMetrics,
			expectedErr:  nil,
		},
		{
			desc:         "No Prefix, minute, logs",
			entityName:   "year=2023/month=01/day=04/hour=12/minute=02/entitylogs_12345.json",
			expectedTime: &expectedTimeMinute,
			expectedType: component.DataTypeLogs,
			expectedErr:  nil,
		},
		{
			desc:         "No Prefix, minute, traces",
			entityName:   "year=2023/month=01/day=04/hour=12/minute=02/entitytraces_12345.json",
			expectedTime: &expectedTimeMinute,
			expectedType: component.DataTypeTraces,
			expectedErr:  nil,
		},
		{
			desc:         "No Prefix, hour, metrics",
			entityName:   "year=2023/month=01/day=04/hour=12/entitymetrics_12345.json",
			expectedTime: &expectedTimeHour,
			expectedType: component.DataTypeMetrics,
			expectedErr:  nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			actualTime, actualType, err := ParseEntityPath(tc.entityName)
			if tc.expectedErr != nil {
				require.ErrorContains(t, err, tc.expectedErr.Error())
				require.Nil(t, tc.expectedTime)
			} else {
				require.NoError(t, err)
				require.NotNil(t, actualTime)
				require.True(t, tc.expectedTime.Equal(*actualTime))
				require.Equal(t, tc.expectedType, actualType)
			}
		})
	}
}
