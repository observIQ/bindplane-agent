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

package azureblobrehydrationreceiver

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	testCases := []struct {
		desc      string
		cfg       *Config
		expectErr error
	}{
		{
			desc: "Missing connection string",
			cfg: &Config{
				ConnectionString: "",
				Container:        "container",
				RootFolder:       "root",
				StartingTime:     "2023-10-02T17:00",
				EndingTime:       "2023-10-02T17:01",
				DeleteOnRead:     false,
			},
			expectErr: errors.New("connection_string is required"),
		},
		{
			desc: "Missing container",
			cfg: &Config{
				ConnectionString: "connection_string",
				Container:        "",
				RootFolder:       "root",
				StartingTime:     "2023-10-02T17:00",
				EndingTime:       "2023-10-02T17:01",
				DeleteOnRead:     false,
			},
			expectErr: errors.New("container is required"),
		},
		{
			desc: "Missing starting_time",
			cfg: &Config{
				ConnectionString: "connection_string",
				Container:        "container",
				RootFolder:       "root",
				StartingTime:     "",
				EndingTime:       "2023-10-02T17:01",
				DeleteOnRead:     false,
			},
			expectErr: errors.New("starting_time is invalid: missing value"),
		},
		{
			desc: "Missing ending_time",
			cfg: &Config{
				ConnectionString: "connection_string",
				Container:        "container",
				RootFolder:       "root",
				StartingTime:     "2023-10-02T17:00",
				EndingTime:       "",
				DeleteOnRead:     false,
			},
			expectErr: errors.New("ending_time is invalid: missing value"),
		},
		{
			desc: "Invalid starting_time",
			cfg: &Config{
				ConnectionString: "connection_string",
				Container:        "container",
				RootFolder:       "root",
				StartingTime:     "invalid_time",
				EndingTime:       "2023-10-02T17:01",
				DeleteOnRead:     false,
			},
			expectErr: errors.New("starting_time is invalid: invalid timestamp"),
		},
		{
			desc: "Missing ending_time",
			cfg: &Config{
				ConnectionString: "connection_string",
				Container:        "container",
				RootFolder:       "root",
				StartingTime:     "2023-10-02T17:00",
				EndingTime:       "invalid_time",
				DeleteOnRead:     false,
			},
			expectErr: errors.New("ending_time is invalid: invalid timestamp"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.expectErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectErr.Error())
			}
		})
	}
}
