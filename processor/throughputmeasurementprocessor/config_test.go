package throughputmeasurementprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		cfg         Config
		expectedErr error
	}{
		{
			desc: "Not enabled",
			cfg: Config{
				Enabled:       false,
				SamplingRatio: 1.0,
			},
			expectedErr: nil,
		},
		{
			desc: "Bad sampling ratio",
			cfg: Config{
				Enabled:       true,
				SamplingRatio: 2.0,
			},
			expectedErr: errInvalidSamplingRatio,
		},
		{
			desc: "Valid config",
			cfg: Config{
				Enabled:       true,
				SamplingRatio: 0.5,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualErr := tc.cfg.Validate()
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}
