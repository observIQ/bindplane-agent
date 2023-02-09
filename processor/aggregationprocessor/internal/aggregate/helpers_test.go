package aggregate

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestGetDatapointValueDouble(t *testing.T) {
	t.Run("with int val", func(t *testing.T) {
		dp := pmetric.NewNumberDataPoint()
		dp.SetIntValue(10)
		require.Equal(t, float64(10), getDatapointValueDouble(dp))
	})

	t.Run("with double val", func(t *testing.T) {
		dp := pmetric.NewNumberDataPoint()
		dp.SetDoubleValue(14.5)
		require.Equal(t, float64(14.5), getDatapointValueDouble(dp))
	})

	t.Run("with no val", func(t *testing.T) {
		dp := pmetric.NewNumberDataPoint()
		require.Equal(t, float64(0), getDatapointValueDouble(dp))
	})
}

func TestGetDatapointValueInt(t *testing.T) {
	t.Run("with int val", func(t *testing.T) {
		dp := pmetric.NewNumberDataPoint()
		dp.SetIntValue(10)
		require.Equal(t, int64(10), getDatapointValueInt(dp))
	})

	t.Run("with double val", func(t *testing.T) {
		dp := pmetric.NewNumberDataPoint()
		dp.SetDoubleValue(14.5)
		require.Equal(t, int64(14), getDatapointValueInt(dp))
	})

	t.Run("with no val", func(t *testing.T) {
		dp := pmetric.NewNumberDataPoint()
		require.Equal(t, int64(0), getDatapointValueInt(dp))
	})
}
