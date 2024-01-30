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
				Database:          "database",
				Warehouse:         "warehouse",
				Metrics: &TelemetryConfig{
					Schema: "schema",
					Table:  "table",
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
				Metrics: &TelemetryConfig{
					Schema: "",
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
				Metrics: &TelemetryConfig{
					Schema: "",
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
				Metrics: &TelemetryConfig{
					Schema: "",
				},
			},
			expectedErr: errors.New("password is required"),
		},
		{
			desc: "Missing database",
			cfg: &Config{
				AccountIdentifier: "id",
				Username:          "user",
				Password:          "pass",
				Warehouse:         "wh",
				Metrics: &TelemetryConfig{
					Schema: "",
				},
			},
			expectedErr: errors.New("database is required"),
		},
		{
			desc: "Missing warehouse",
			cfg: &Config{
				AccountIdentifier: "id",
				Username:          "user",
				Password:          "pass",
				Database:          "db",
				Metrics: &TelemetryConfig{
					Schema: "",
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
				Logs: &TelemetryConfig{
					Schema: "",
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
				Metrics: &TelemetryConfig{
					Schema: "",
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
				Traces: &TelemetryConfig{
					Schema: "",
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
