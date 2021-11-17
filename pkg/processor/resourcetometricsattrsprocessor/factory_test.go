package resourcetometricsattrsprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configtest"
	"go.opentelemetry.io/collector/consumer/consumertest"
)

func TestNewFactory(t *testing.T) {
	f := NewFactory()
	require.NotNil(t, f)
}

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	require.NotNil(t, cfg)
	require.NoError(t, configtest.CheckConfigStruct(cfg))
}

func TestCreateMetricsExporter(t *testing.T) {
	cfg := createDefaultConfig()
	p, err := createMetricsProcessor(context.Background(), componenttest.NewNopProcessorCreateSettings(), cfg, consumertest.NewNop())
	require.NotNil(t, p)
	require.NoError(t, err)
}

func TestCreateMetricsExporterNilConfig(t *testing.T) {
	_, err := createMetricsProcessor(context.Background(), componenttest.NewNopProcessorCreateSettings(), nil, consumertest.NewNop())
	require.Error(t, err)
}
