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

package githubreceiver // import "github.com/observiq/bindplane-agent/receiver/githubreceiver"

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		errExpected bool
		errText     string
		config      Config
	}{
		{
			desc:        "pass simple",
			errExpected: false,
			config: Config{
				AccessToken:  "AccessToken",
				LogType:      "user",
				Name:         "testName",
				PollInterval: time.Second,
			},
		},
		{
			desc:        "fail no access token",
			errExpected: true,
			errText:     "missing access_token; required",
			config: Config{
				LogType:      "user",
				Name:         "testName",
				PollInterval: time.Second,
			},
		},
		{
			desc:        "fail no log type",
			errExpected: true,
			errText:     "missing log_type; required",
			config: Config{
				AccessToken:  "AccessToken",
				Name:         "testName",
				PollInterval: time.Second,
			},
		},
		{
			desc:        "fail no name",
			errExpected: true,
			errText:     "missing name; required",
			config: Config{
				AccessToken:  "AccessToken",
				LogType:      "user",
				PollInterval: time.Second,
			},
		},
		{
			desc:        "fail no name",
			errExpected: true,
			errText:     "missing name; required",
			config: Config{
				AccessToken:  "AccessToken",
				LogType:      "user",
				PollInterval: time.Second,
			},
		},
		{
			desc:        "fail with no poll interval or no webhook",
			errExpected: true,
			errText:     "must specify either poll_interval or webhook",
			config: Config{
				AccessToken: "AccessToken",
				LogType:     "user",
				Name:        "testName",
			},
		},
		{
			desc:        "fail invalid poll interval short",
			errExpected: true,
			errText:     "invalid poll_interval; must be at least 0.72 seconds",
			config: Config{
				AccessToken:  "AccessToken",
				LogType:      "user",
				Name:         "testName",
				PollInterval: 700 * time.Millisecond,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			err := component.ValidateConfig(tc.config)

			if tc.errExpected {
				require.EqualError(t, err, tc.errText)
				return
			}

			require.NoError(t, err)
		})
	}
}
