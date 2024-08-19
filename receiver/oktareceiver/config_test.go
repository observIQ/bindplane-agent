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

package oktareceiver

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
				Domain:   "oktadomain.com",
				ApiToken: "dummyApiToken",
			},
		},
		{
			desc: "pass with poll interval",
			config: Config{
				Domain:       "oktadomain.com",
				ApiToken:     "dummyApiToken",
				PollInterval: time.Second,
			},
		},
		{
			desc: "pass with poll interval zero value",
			config: Config{
				Domain:       "oktadomain.com",
				ApiToken:     "dummyApiToken",
				PollInterval: 0,
			},
		},
		{
			desc:        "fail no domain",
			expectedErr: errNoDomain,
			config: Config{
				ApiToken: "dummyApiToken",
			},
		},
		{
			desc:        "fail no api token",
			expectedErr: errNoApiToken,
			config: Config{
				Domain: "oktadomain.com",
			},
		},

		{
			desc:        "fail invalid poll interval",
			expectedErr: errInvalidPollInterval,
			config: Config{
				Domain:       "oktadomain.com",
				ApiToken:     "dummyApiToken",
				PollInterval: 500 * time.Millisecond,
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
