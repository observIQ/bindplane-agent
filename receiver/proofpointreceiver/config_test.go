// Copyright observIQ, Inc.
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

package proofpointreceiver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		expectedErr error
		config      Config
	}{
		{
			desc: "pass simple",
			config: Config{
				Principal:    "dummyPrincipal",
				Secret:       "dummySecret",
				PollInterval: 5 * time.Minute,
			},
		},
		{
			desc:        "fail no principal",
			expectedErr: errNoPrincipal,
			config: Config{
				Secret:       "dummySecret",
				PollInterval: 5 * time.Minute,
			},
		},
		{
			desc:        "fail no secret",
			expectedErr: errNoSecret,
			config: Config{
				Principal:    "dummyPrincipal",
				PollInterval: 5 * time.Minute,
			},
		},
		{
			desc:        "fail invalid poll_interval short",
			expectedErr: errInvalidPollInterval,
			config: Config{
				Principal:    "dummyPrincipal",
				Secret:       "dummySecret",
				PollInterval: 59 * time.Second,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := component.ValidateConfig(tc.config)
			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
