package pluginreceiver

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestWrapLogger(t *testing.T) {
	baseLogger := zap.NewNop()
	opts := createServiceLoggerOpts(baseLogger)
	serviceLogger := zap.NewNop().WithOptions(opts...)
	require.Equal(t, baseLogger.Core(), serviceLogger.Core())

	infoLevel := serviceLogger.Core().Enabled(zapcore.InfoLevel)
	require.False(t, infoLevel)
}

func TestCreateService(t *testing.T) {
	configProvider := createConfigProvider(nil)
	factories := component.Factories{}
	logger := zap.NewNop()
	_, err := createService(factories, configProvider, logger)
	require.NoError(t, err)
}
