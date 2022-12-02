package logcountprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

func TestNewProcessorFactory(t *testing.T) {
	f := NewProcessorFactory()
	require.Equal(t, component.NewID(typeStr).Type(), f.Type())
	require.Equal(t, stability, f.LogsProcessorStability())
	require.NotNil(t, f.CreateDefaultConfig())
	require.NotNil(t, f.CreateLogsProcessor)
}

func TestNewReceiverFactory(t *testing.T) {
	f := NewReceiverFactory()
	require.Equal(t, component.NewID(typeStr).Type(), f.Type())
	require.Equal(t, stability, f.MetricsReceiverStability())
	require.NotNil(t, f.CreateDefaultConfig())
	require.NotNil(t, f.CreateMetricsReceiver)
}

func TestCreateLogsProcessor(t *testing.T) {
	var testCases = []struct {
		name        string
		cfg         component.ProcessorConfig
		expectedErr string
	}{
		{
			name: "valid config",
			cfg: &ProcessorConfig{
				ProcessorSettings: config.ProcessorSettings{},
				Match:             "true",
			},
		},
		{
			name: "invalid match",
			cfg: &ProcessorConfig{
				ProcessorSettings: config.ProcessorSettings{},
				Match:             "++",
			},
			expectedErr: "failed to create match expression",
		},
		{
			name: "invalid attribute",
			cfg: &ProcessorConfig{
				ProcessorSettings: config.ProcessorSettings{},
				Match:             "true",
				Attributes:        map[string]string{"a": "++"},
			},
			expectedErr: "failed to create attribute expression for a",
		},
		{
			name:        "invalid config type",
			cfg:         &ReceiverConfig{},
			expectedErr: "invalid config type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewProcessorFactory()
			p, err := f.CreateLogsProcessor(context.Background(), component.ProcessorCreateSettings{}, tc.cfg, nil)
			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.IsType(t, &processor{}, p)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			}
		})
	}
}

func TestCreateMetricsReceiver(t *testing.T) {
	f := NewReceiverFactory()
	r, err := f.CreateMetricsReceiver(context.Background(), component.ReceiverCreateSettings{}, createDefaultReceiverConfig(), nil)
	require.NoError(t, err)
	require.IsType(t, &receiver{}, r)
}
