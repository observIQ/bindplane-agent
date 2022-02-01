package collector

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewSettings(t *testing.T) {
	settings := NewSettings("./test/valid.yaml", "0.0.0", nil)
	require.Equal(t, settings.LoggingOptions, []zap.Option(nil))
	require.True(t, settings.DisableGracefulShutdown)
}
