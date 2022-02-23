package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name      string
		input     Config
		expectErr bool
	}{
		{
			"missing-namespace",
			Config{
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
			Config{
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
			Config{
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
			Config{
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
			Config{
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
			Config{
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
			Config{
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
			Config{
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
