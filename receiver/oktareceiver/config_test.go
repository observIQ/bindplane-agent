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
	"errors"
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
				APIToken: "dummyAPIToken",
			},
		},
		{
			desc: "pass with poll interval",
			config: Config{
				Domain:       "oktadomain.com",
				APIToken:     "dummyAPIToken",
				PollInterval: time.Second,
			},
		},
		{
			desc: "pass with start time",
			config: Config{
				Domain:    "oktadomain.com",
				APIToken:  "dummyAPIToken",
				StartTime: time.Now().Add(-time.Hour).Format(OktaTimeFormat),
			},
		},
		{
			desc: "pass with all fields",
			config: Config{
				Domain:       "oktadomain.com",
				APIToken:     "dummyAPIToken",
				PollInterval: time.Second,
				StartTime:    time.Now().Add(-time.Hour).Format(OktaTimeFormat),
			},
		},
		{
			desc: "pass with poll interval zero value",
			config: Config{
				Domain:       "oktadomain.com",
				APIToken:     "dummyAPIToken",
				PollInterval: 0,
			},
		},
		{
			desc:        "fail no domain",
			expectedErr: errNoDomain,
			config: Config{
				APIToken: "dummyAPIToken",
			},
		},
		{
			desc:        "fail invalid domain https",
			expectedErr: errInvalidDomain,
			config: Config{
				APIToken: "dummyAPIToken",
				Domain:   "https://test.okta.com",
			},
		},
		{
			desc:        "fail invalid domain http",
			expectedErr: errInvalidDomain,
			config: Config{
				APIToken: "dummyAPIToken",
				Domain:   "http://test.okta.com",
			},
		},
		{
			desc:        "fail no api token",
			expectedErr: errNoAPIToken,
			config: Config{
				Domain: "oktadomain.com",
			},
		},
		{
			desc:        "fail invalid poll interval",
			expectedErr: errInvalidPollInterval,
			config: Config{
				Domain:       "oktadomain.com",
				APIToken:     "dummyAPIToken",
				PollInterval: 500 * time.Millisecond,
			},
		},
		{
			desc:        "fail invalid start time format",
			expectedErr: errors.New("invalid start_time: invalid timestamp: must be in the format YYYY-MM-DDTHH:MM:SS"),
			config: Config{
				Domain:    "oktadomain.com",
				APIToken:  "dummyAPIToken",
				StartTime: time.Now().UTC().Add(-time.Hour).Format(time.RFC1123),
			},
		},
		{
			desc:        "fail invalid start time future",
			expectedErr: errors.New("invalid start_time: invalid timestamp: must be within the past 180 days and not in the future"),
			config: Config{
				Domain:    "oktadomain.com",
				APIToken:  "dummyAPIToken",
				StartTime: time.Now().UTC().Add(time.Hour).Format(OktaTimeFormat),
			},
		},
		{
			desc:        "fail invalid start time too old",
			expectedErr: errors.New("invalid start_time: invalid timestamp: must be within the past 180 days and not in the future"),
			config: Config{
				Domain:    "oktadomain.com",
				APIToken:  "dummyAPIToken",
				StartTime: time.Now().UTC().AddDate(0, 0, -181).Format(OktaTimeFormat),
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
