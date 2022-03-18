package pluginreceiver

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

func TestCreateReceiver(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config.Receiver
		expectedErr error
	}{
		{
			name:        "invalid config type",
			cfg:         &config.ReceiverSettings{},
			expectedErr: errors.New("config is not a plugin receiver config"),
		},
		{
			name: "missing plugin",
			cfg: &Config{
				Path: "./test/missing.yaml",
			},
			expectedErr: errors.New("failed to load plugin"),
		},
		{
			name: "invalid plugin yaml",
			cfg: &Config{
				Path: "./test/plugin-invalid-yaml.yaml",
			},
			expectedErr: errors.New("failed to load plugin"),
		},
		{
			name: "invalid plugin parameter",
			cfg: &Config{
				Path: "./test/plugin-valid.yaml",
				Parameters: map[string]interface{}{
					"env": 5,
				},
			},
			expectedErr: errors.New("invalid plugin parameter"),
		},
		{
			name: "invalid plugin template",
			cfg: &Config{
				Path: "./test/plugin-invalid-template.yaml",
				Parameters: map[string]interface{}{
					"env": "prod",
				},
			},
			expectedErr: errors.New("failed to render plugin"),
		},
		{
			name: "valid plugin",
			cfg: &Config{
				Path: "./test/plugin-valid.yaml",
				Parameters: map[string]interface{}{
					"env": "prod",
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			set := component.ReceiverCreateSettings{}
			consumer := &MockConsumer{}
			emitterFactory := createLogEmitterFactory(consumer)
			receiver, err := createReceiver(tc.cfg, set, emitterFactory, config.LogsDataType)

			switch tc.expectedErr {
			case nil:
				require.NoError(t, err)
				require.IsType(t, &Receiver{}, receiver)
			default:
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr.Error())
			}
		})
	}
}

func TestCreateLogsReceiver(t *testing.T) {
	factory := NewFactory()
	consumer := &MockConsumer{}
	ctx := context.Background()
	set := component.ReceiverCreateSettings{}
	cfg := &Config{
		Path: "./test/plugin-valid.yaml",
		Parameters: map[string]interface{}{
			"env": "prod",
		},
	}

	receiver, err := factory.CreateLogsReceiver(ctx, set, cfg, consumer)
	require.NoError(t, err)
	require.IsType(t, &Receiver{}, receiver)
}

func TestCreateMetricsReceiver(t *testing.T) {
	factory := NewFactory()
	consumer := &MockConsumer{}
	ctx := context.Background()
	set := component.ReceiverCreateSettings{}
	cfg := &Config{
		Path: "./test/plugin-valid.yaml",
		Parameters: map[string]interface{}{
			"env": "prod",
		},
	}

	receiver, err := factory.CreateMetricsReceiver(ctx, set, cfg, consumer)
	require.NoError(t, err)
	require.IsType(t, &Receiver{}, receiver)
}

func TestCreateTracesReceiver(t *testing.T) {
	factory := NewFactory()
	consumer := &MockConsumer{}
	ctx := context.Background()
	set := component.ReceiverCreateSettings{}
	cfg := &Config{
		Path: "./test/plugin-valid.yaml",
		Parameters: map[string]interface{}{
			"env": "prod",
		},
	}

	receiver, err := factory.CreateTracesReceiver(ctx, set, cfg, consumer)
	require.NoError(t, err)
	require.IsType(t, &Receiver{}, receiver)
}

func TestCreateDefaultConfig(t *testing.T) {
	config := createDefaultConfig()
	require.Equal(t, typeStr, config.ID().String())

	pluginConfig, ok := config.(*Config)
	require.True(t, ok)
	require.Equal(t, make(map[string]interface{}), pluginConfig.Parameters)
	require.Empty(t, pluginConfig.Path)
}
