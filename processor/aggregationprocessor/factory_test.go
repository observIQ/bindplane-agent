package aggregationprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestNewFactory(t *testing.T) {
	t.Run("default config is valid", func(t *testing.T) {
		fact := NewFactory()
		require.NotNil(t, fact)

		conf := fact.CreateDefaultConfig()

		c, ok := conf.(*Config)
		require.True(t, ok, "default config from factory was not processor config!")

		// Default config should be valid
		require.NoError(t, c.Validate())
	})

	t.Run("metrics processor is created with default config", func(t *testing.T) {
		fact := NewFactory()
		require.NotNil(t, fact)

		conf := fact.CreateDefaultConfig()

		c, ok := conf.(*Config)
		require.True(t, ok, "default config from factory was not processor config!")

		// Default config should be valid
		require.NoError(t, c.Validate())

		defaultProcessor, err := fact.CreateMetricsProcessor(
			context.Background(),
			processortest.NewNopCreateSettings(),
			conf,
			consumertest.NewNop(),
		)

		require.NoError(t, err)
		require.NotNil(t, defaultProcessor)
	})

	t.Run("metrics processor fails to create with incorrect config type", func(t *testing.T) {
		fact := NewFactory()
		require.NotNil(t, fact)

		_, err := fact.CreateMetricsProcessor(
			context.Background(),
			processortest.NewNopCreateSettings(),
			"not a config",
			consumertest.NewNop(),
		)

		require.ErrorContains(t, err, "cannot create aggregation processor with invalid config type:")
	})

	t.Run("metrics processor fails to create with bad regex", func(t *testing.T) {
		fact := NewFactory()
		require.NotNil(t, fact)

		_, err := fact.CreateMetricsProcessor(
			context.Background(),
			processortest.NewNopCreateSettings(),
			&Config{
				Include: "^(",
			},
			consumertest.NewNop(),
		)

		require.ErrorContains(t, err, "failed to create aggregation processor:")
	})

}
