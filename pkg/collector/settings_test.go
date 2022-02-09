package collector

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewSettings(t *testing.T) {
	configPaths := []string{"./test/valid.yaml"}
	settings := NewSettings(configPaths, "0.0.0", nil)
	require.Equal(t, settings.LoggingOptions, []zap.Option(nil))
	require.True(t, settings.DisableGracefulShutdown)
}
