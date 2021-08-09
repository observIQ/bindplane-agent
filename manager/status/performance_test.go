package status

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetricCollection(t *testing.T) {
	sr := Report{
		ComponentType: "observiq-collector",
		ComponentID:   "id",
		Status:        Status(ACTIVE),
		Metrics:       map[string]*Metric{},
	}
	err := sr.addPerformanceMetrics()
	require.NoError(t, err)
	if value, hasValue := sr.Metrics[string(CPU_PERCENT)]; hasValue {
		v, isFloat := value.Value.(float64)
		require.True(t, isFloat)
		require.GreaterOrEqual(t, v, 0.0)
	}

	if val, hv := sr.Metrics[string(MEMORY_USED)]; hv {
		v, isFloat := val.Value.(float64)
		require.True(t, isFloat)
		require.GreaterOrEqual(t, v, 0.0)
	}

	if inVal, hasNetworkIn := sr.Metrics[string(NETWORK_DATA_IN)]; hasNetworkIn {
		_, isFloat := inVal.Value.(uint64)
		require.True(t, isFloat)
	}

	if inVal, hasNetworkOut := sr.Metrics[string(NETWORK_DATA_OUT)]; hasNetworkOut {
		_, isFloat := inVal.Value.(uint64)
		require.True(t, isFloat)
	}
}
