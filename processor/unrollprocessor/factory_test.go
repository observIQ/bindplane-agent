package unrollprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	require.Equal(t, componentType, factory.Type())

	expectedCfg := &Config{
		Field: UnrollFieldBody,
	}

	cfg, ok := factory.CreateDefaultConfig().(*Config)
	require.True(t, ok)
	require.Equal(t, expectedCfg, cfg)
}

func TestBadFactory(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)
	cfg.Field = "invalid"

	_, err := factory.CreateLogs(context.Background(), processortest.NewNopSettings(), cfg, &consumertest.LogsSink{})
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid config for \"unroll\" processor")
}
