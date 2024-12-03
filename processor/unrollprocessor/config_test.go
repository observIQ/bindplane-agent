package unrollprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {

	testCases := []struct {
		desc        string
		cfg         *Config
		expectedErr string
	}{
		{
			desc: "valid config",
			cfg:  createDefaultConfig().(*Config),
		},
		{
			desc: "config without body field",
			cfg: &Config{
				Field: "attributes",
			},
			expectedErr: "only unrolling logs from a body slice is currently supported",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
