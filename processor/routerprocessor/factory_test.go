package routerprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

func TestNewFactory(t *testing.T) {
	expectedCfg := &Config{
		ProcessorSettings: config.NewProcessorSettings(component.NewID(typeStr)),
		Routes:            make([]*RoutingRule, 0),
	}

	f := NewFactory()
	require.Equal(t, component.Type(typeStr), f.Type())
	defaultCfg := f.CreateDefaultConfig()
	require.Equal(t, expectedCfg, defaultCfg)
}

func TestCreateLogsProcessor(t *testing.T) {
	var testCases = []struct {
		name        string
		cfg         component.Config
		expectedErr string
	}{
		{
			name: "valid config",
			cfg: &Config{
				ProcessorSettings: config.ProcessorSettings{},
				Routes: []*RoutingRule{
					{
						Match: "true",
						Route: "route1",
					},
				},
			},
		},
		{
			name: "invalid match",
			cfg: &Config{
				ProcessorSettings: config.ProcessorSettings{},
				Routes: []*RoutingRule{
					{
						Match: "++",
						Route: "route1",
					},
				},
			},
			expectedErr: "invalid match expression",
		},
		{
			name:        "invalid config type",
			cfg:         nil,
			expectedErr: "invalid config type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewFactory()
			p, err := f.CreateLogsProcessor(context.Background(), component.ProcessorCreateSettings{}, tc.cfg, nil)
			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.IsType(t, &routerProcessor{}, p)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			}
		})
	}
}
