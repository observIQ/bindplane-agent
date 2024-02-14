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

package snowflakeexporter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		cfg         *Config
		expectedErr error
	}{
		{
			desc: "Simple metrics pass",
			cfg: &Config{
				AccountIdentifier: "accountID",
				Username:          "username",
				Password:          "password",
				Warehouse:         "warehouse",
				Metrics: TelemetryConfig{
					Enabled: true,
				},
			},
		},
		{
			desc: "Missing account identifier",
			cfg: &Config{
				Username:  "user",
				Password:  "pass",
				Database:  "db",
				Warehouse: "wh",
				Metrics: TelemetryConfig{
					Enabled: true,
				},
			},
			expectedErr: errors.New("account_identifier is required"),
		},
		{
			desc: "Missing username",
			cfg: &Config{
				AccountIdentifier: "id",
				Password:          "pass",
				Database:          "db",
				Warehouse:         "wh",
				Metrics: TelemetryConfig{
					Enabled: true,
				},
			},
			expectedErr: errors.New("username is required"),
		},
		{
			desc: "Missing password",
			cfg: &Config{
				AccountIdentifier: "id",
				Username:          "user",
				Database:          "db",
				Warehouse:         "wh",
				Metrics: TelemetryConfig{
					Enabled: true,
				},
			},
			expectedErr: errors.New("password is required"),
		},
		{
			desc: "Missing warehouse",
			cfg: &Config{
				AccountIdentifier: "id",
				Username:          "user",
				Password:          "pass",
				Database:          "db",
				Metrics: TelemetryConfig{
					Enabled: true,
				},
			},
			expectedErr: errors.New("warehouse is required"),
		},
		{
			desc: "Default logs cfg",
			cfg: &Config{
				AccountIdentifier: "id",
				Username:          "user",
				Password:          "pass",
				Database:          "db",
				Warehouse:         "wh",
				Logs: TelemetryConfig{
					Enabled: true,
				},
			},
		},
		{
			desc: "Default metrics cfg",
			cfg: &Config{
				AccountIdentifier: "id",
				Username:          "user",
				Password:          "pass",
				Database:          "db",
				Warehouse:         "wh",
				Metrics: TelemetryConfig{
					Enabled: true,
				},
			},
		},
		{
			desc: "Default traces cfg",
			cfg: &Config{
				AccountIdentifier: "id",
				Username:          "user",
				Password:          "pass",
				Database:          "db",
				Warehouse:         "wh",
				Traces: TelemetryConfig{
					Enabled: true,
				},
			},
		},
		{
			desc: "Partial telemetry cfgs",
			cfg: &Config{
				AccountIdentifier: "id",
				Username:          "user",
				Password:          "pass",
				Warehouse:         "wh",
				Role:              "role",
				Logs: TelemetryConfig{
					Enabled: true,
					Table:   "lt",
				},
				Metrics: TelemetryConfig{
					Enabled: true,
					Schema:  "ms",
				},
				Traces: TelemetryConfig{
					Enabled: true,
				},
			},
		},
		{
			desc: "No telemetry configured",
			cfg: &Config{
				AccountIdentifier: "id",
				Username:          "user",
				Password:          "pass",
				Database:          "db",
				Warehouse:         "wh",
			},
			expectedErr: errors.New("no telemetry type configured for exporter"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}
