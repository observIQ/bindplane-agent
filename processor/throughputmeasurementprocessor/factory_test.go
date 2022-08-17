package throughputmeasurementprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config"
)

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	require.Equal(t, typeStr, string(factory.Type()))

	expectedCfg := &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
		Enabled:           true,
		SamplingRatio:     1.0,
	}

	cfg, ok := factory.CreateDefaultConfig().(*Config)
	require.True(t, ok)
	require.Equal(t, expectedCfg, cfg)
}
