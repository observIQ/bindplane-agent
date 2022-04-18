package googlecloudexporter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config"
)

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	googleCfg, ok := cfg.(*Config)
	require.True(t, ok)

	require.Equal(t, config.NewComponentID("googlecloud"), googleCfg.ID())
	require.Equal(t, defaultMetricPrefix, googleCfg.GCPConfig.MetricConfig.Prefix)
	require.Equal(t, defaultUserAgent, googleCfg.GCPConfig.UserAgent)
	require.Nil(t, googleCfg.Validate())
}
