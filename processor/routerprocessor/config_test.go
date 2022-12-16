package routerprocessor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

func TestConfigValidate(t *testing.T) {
	testCases := []struct {
		desc        string
		cfg         Config
		expectedErr error
	}{
		{
			desc: "No routes specified",
			cfg: Config{
				ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
			},
			expectedErr: errNoRoutesSpecified,
		},
		{
			desc: "Duplicate match statements",
			cfg: Config{
				ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
				Routes: []*RoutingRule{
					{
						Match: "match1",
						Route: "route1",
					},
					{
						Match: "match1",
						Route: "route2",
					},
				},
			},
			expectedErr: errors.New("duplicate match expression"),
		},
		{
			desc: "Duplicate route name",
			cfg: Config{
				ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
				Routes: []*RoutingRule{
					{
						Match: "match1",
						Route: "route1",
					},
					{
						Match: "match2",
						Route: "route1",
					},
				},
			},
			expectedErr: errors.New("duplicate route name"),
		},
		{
			desc: "Valid config",
			cfg: Config{
				ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
				Routes: []*RoutingRule{
					{
						Match: "match1",
						Route: "route1",
					},
					{
						Match: "match2",
						Route: "route2",
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualErr := tc.cfg.Validate()
			if tc.expectedErr != nil {
				require.ErrorContains(t, actualErr, tc.expectedErr.Error())
			} else {
				require.NoError(t, actualErr)
			}
		})
	}
}
