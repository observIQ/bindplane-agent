package logcountprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
)

func TestCreateDefaultProcessorConfig(t *testing.T) {
	cfg := createDefaultProcessorConfig().(*ProcessorConfig)
	require.Equal(t, defaultInterval, cfg.Interval)
	require.Equal(t, defaultMatch, cfg.Match)
	require.Equal(t, defaultMetricName, cfg.MetricName)
	require.Equal(t, defaultMetricUnit, cfg.MetricUnit)
	require.Equal(t, component.NewID(typeStr), cfg.ProcessorSettings.ID())
}

func TestCreateDefaultReceiverConfig(t *testing.T) {
	cfg := createDefaultReceiverConfig().(*ReceiverConfig)
	require.Equal(t, component.NewID(typeStr), cfg.ReceiverSettings.ID())
}

func TestCreateMatchExpr(t *testing.T) {
	cfg := createDefaultProcessorConfig().(*ProcessorConfig)
	cfg.Match = "true"
	expr, err := cfg.createMatchExpr()
	require.NoError(t, err)
	require.NotNil(t, expr)

	cfg.Match = "++"
	expr, err = cfg.createMatchExpr()
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create match expression")
}

func TestCreateAttrExprs(t *testing.T) {
	cfg := createDefaultProcessorConfig().(*ProcessorConfig)
	cfg.Attributes = map[string]string{"a": "true"}
	expr, err := cfg.createAttrExprs()
	require.NoError(t, err)
	require.NotNil(t, expr)

	cfg.Attributes = map[string]string{"a": "++"}
	expr, err = cfg.createAttrExprs()
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to create attribute expression for a")
}
