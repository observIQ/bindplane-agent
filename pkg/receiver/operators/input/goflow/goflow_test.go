package goflow

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-log-collection/testutil"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	cases := []struct {
		name        string
		inputRecord InputConfig
		expectErr   bool
	}{
		{
			"minimal-default-mode",
			InputConfig{
				ListenAddress: "0.0.0.0:2056",
			},
			false,
		},
		{
			"minimal-netflow-v5",
			InputConfig{
				Mode:          "netflow_v5",
				ListenAddress: "0.0.0.0:2056",
			},
			false,
		},
		{
			"minimal-netflow-ipfix",
			InputConfig{
				Mode:          "netflow_ipfix",
				ListenAddress: "0.0.0.0:2056",
			},
			false,
		},
		{
			"minimal-netflow-sflow",
			InputConfig{
				Mode:          "netflow_v5",
				ListenAddress: "0.0.0.0:2056",
			},
			false,
		},
		{
			"invalid mode",
			InputConfig{
				Mode:          "netflow",
				ListenAddress: "0.0.0.0:2056",
			},
			true,
		},
		{
			"missing-address",
			InputConfig{
				Mode: "sflow",
			},
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewGoflowInputConfig("test_id")
			cfg.ListenAddress = tc.inputRecord.ListenAddress
			cfg.Mode = tc.inputRecord.Mode

			_, err := cfg.Build(testutil.NewBuildContext(t))
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
