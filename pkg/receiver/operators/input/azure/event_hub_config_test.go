package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name      string
		input     AzureConfig
		expectErr bool
	}{
		{
			"missing-namespace",
			AzureConfig{
				Namespace:        "",
				Name:             "john",
				Group:            "devel",
				ConnectionString: "some connection string",
				StartAt:          "end",
				PrefetchCount:    10,
			},
			true,
		},
		{
			"missing-name",
			AzureConfig{
				Namespace:        "namespace",
				Name:             "",
				Group:            "devel",
				ConnectionString: "some connection string",
				StartAt:          "end",
				PrefetchCount:    10,
			},
			true,
		},
		{
			"missing-group",
			AzureConfig{
				Namespace:        "namespace",
				Name:             "dev",
				Group:            "",
				ConnectionString: "some connection string",
				StartAt:          "end",
				PrefetchCount:    10,
			},
			true,
		},
		{
			"missing-connection-string",
			AzureConfig{
				Namespace:        "namespace",
				Name:             "dev",
				Group:            "dev",
				ConnectionString: "",
				StartAt:          "end",
				PrefetchCount:    10,
			},
			true,
		},
		{
			"invalid-prefetch-count",
			AzureConfig{
				Namespace:        "namespace",
				Name:             "dev",
				Group:            "dev",
				ConnectionString: "some string",
				StartAt:          "end",
				PrefetchCount:    0,
			},
			true,
		},
		{
			"invalid-start-at",
			AzureConfig{
				Namespace:        "namespace",
				Name:             "dev",
				Group:            "dev",
				ConnectionString: "some string",
				StartAt:          "bad",
				PrefetchCount:    10,
			},
			true,
		},
		{
			"valid-start-at-end",
			AzureConfig{
				Namespace:        "namespace",
				Name:             "dev",
				Group:            "dev",
				ConnectionString: "some string",
				StartAt:          "end",
				PrefetchCount:    10,
			},
			false,
		},
		{
			"valid-start-at-beginning",
			AzureConfig{
				Namespace:        "namespace",
				Name:             "dev",
				Group:            "dev",
				ConnectionString: "some string",
				PrefetchCount:    10,
				StartAt:          "beginning",
			},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.validate()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
