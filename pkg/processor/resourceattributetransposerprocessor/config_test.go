package resourceattributetransposerprocessor

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/service/servicetest"
)

func TestConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	require.NoError(t, err)

	factory := NewFactory()
	factories.Processors[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(path.Join(".", "testdata", "config.yaml"), factories)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	require.Equal(t, len(cfg.Processors), 2)

	// Loaded config should be equal to default config
	defaultCfg := factory.CreateDefaultConfig()
	r0 := cfg.Processors[config.NewComponentID(typeStr)]
	require.Equal(t, r0, defaultCfg)

	customComponentID := config.NewComponentIDWithName(typeStr, "customname")
	r1 := cfg.Processors[customComponentID].(*Config)
	require.Equal(t, &Config{
		ProcessorSettings: config.NewProcessorSettings(customComponentID),
		Operations: []CopyResourceConfig{
			{
				From: "some.resource.level.attr",
				To:   "some.metricdatapoint.level.attr",
			},
			{
				From: "another.resource.attr",
				To:   "another.datapoint.attr",
			},
		},
	}, r1)
}
