package chronicleexporter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		config      *Config
		expectedErr string
	}{
		{
			desc: "Both creds_file_path and creds are empty",
			config: &Config{
				Region:  "United States Multi-Region",
				LogType: "log_type_example",
			},
			expectedErr: "either creds_file_path or creds is required",
		},
		{
			desc: "LogType is empty",
			config: &Config{
				Region: "United States Multi-Region",
				Creds:  "creds_example",
			},
			expectedErr: "log_type is required",
		},
		{
			desc: "Region is empty",
			config: &Config{
				Creds:   "creds_example",
				LogType: "log_type_example",
			},
			expectedErr: "region is required",
		},
		{
			desc: "Region is invalid",
			config: &Config{
				Region:  "Invalid Region",
				Creds:   "creds_example",
				LogType: "log_type_example",
			},
			expectedErr: "region is invalid",
		},
		{
			desc: "Valid config with creds",
			config: &Config{
				Region:  "United States Multi-Region",
				Creds:   "creds_example",
				LogType: "log_type_example",
			},
			expectedErr: "",
		},
		{
			desc: "Valid config with creds_file_path",
			config: &Config{
				Region:        "United States Multi-Region",
				CredsFilePath: "/path/to/creds_file",
				LogType:       "log_type_example",
			},
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			}
		})
	}
}
