package eventhub

import (
	"testing"

	"github.com/observiq/observiq-collector/pkg/receiver/operators/input/azure"
	"github.com/open-telemetry/opentelemetry-log-collection/testutil"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	cases := []struct {
		name      string
		input     InputConfig
		expectErr bool
	}{
		{
			"default",
			InputConfig{
				Config: azure.Config{
					Namespace:        "test",
					Name:             "test",
					Group:            "test",
					ConnectionString: "test",
					PrefetchCount:    1000,
				},
			},
			false,
		},
		{
			"prefetch",
			InputConfig{
				Config: azure.Config{
					Namespace:        "test",
					Name:             "test",
					Group:            "test",
					ConnectionString: "test",
					PrefetchCount:    100,
				},
			},
			false,
		},
		{
			"startat-end",
			InputConfig{
				Config: azure.Config{
					Namespace:        "test",
					Name:             "test",
					Group:            "test",
					ConnectionString: "test",
					StartAt:          "end",
					PrefetchCount:    1000,
				},
			},
			false,
		},
		{
			"startat-beginning",
			InputConfig{
				Config: azure.Config{
					Namespace:        "test",
					Name:             "test",
					Group:            "test",
					ConnectionString: "test",
					StartAt:          "beginning",
					PrefetchCount:    1000,
				},
			},
			false,
		},
		{
			"prefetch-invalid",
			InputConfig{
				Config: azure.Config{
					Namespace:        "test",
					Name:             "test",
					Group:            "test",
					ConnectionString: "test",
					PrefetchCount:    0,
				},
			},
			true,
		},
		{
			"default-required-startat-invalid",
			InputConfig{
				Config: azure.Config{
					Namespace:        "test",
					Name:             "test",
					Group:            "test",
					ConnectionString: "test",
					StartAt:          "invalid",
					PrefetchCount:    1000,
				},
			},
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := NewEventHubConfig("test_id")
			cfg.Namespace = tc.input.Namespace
			cfg.Name = tc.input.Name
			cfg.Group = tc.input.Group
			cfg.ConnectionString = tc.input.ConnectionString

			if tc.input.PrefetchCount != NewEventHubConfig("").PrefetchCount {
				cfg.PrefetchCount = tc.input.PrefetchCount
			}

			if tc.input.StartAt != "" {
				cfg.StartAt = tc.input.StartAt
			}

			_, err := cfg.Build(testutil.NewBuildContext(t))
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
