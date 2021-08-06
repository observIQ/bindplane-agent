package collector

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/service"
	"go.uber.org/zap"
)

func TestSettingsValid(t *testing.T) {
	settings, err := NewSettings("./test/valid.yaml", nil)
	require.NoError(t, err)
	require.Equal(t, settings.LoggingOptions, []zap.Option(nil))
	require.True(t, settings.DisableGracefulShutdown)

	fileProvider, ok := settings.ParserProvider.(*FileProvider)
	require.True(t, ok)
	require.Equal(t, "./test/valid.yaml", fileProvider.filePath)
}

func TestSettingsInvalid(t *testing.T) {
	original := defaultExporters
	defer func() { defaultExporters = original }()
	defaultExporters = append(defaultExporters, loggingexporter.NewFactory())

	settings, err := NewSettings("./test/valid.yaml", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to build factories")
	require.Equal(t, service.CollectorSettings{}, settings)
}
