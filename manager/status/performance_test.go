package status

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetricCollection(t *testing.T) {
	sr := &Report{
		ComponentType: "observiq-collector",
		ComponentID:   "id",
		Status:        ACTIVE,
		Metrics:       map[string]*Metric{},
	}
	sr.AddPerformanceMetrics(nil)

	// CPU metrics are not captured on darwin
	if runtime.GOOS != "darwin" {
		if value, hasValue := sr.Metrics[string(CPU_PERCENT)]; hasValue {
			v, isFloat := value.Value.(float64)
			require.True(t, isFloat)
			require.GreaterOrEqual(t, v, 0.0)
		} else {
			require.FailNow(t, "Did not attach CPU percent metric")
		}
	}

	if val, hv := sr.Metrics[string(MEMORY_USED)]; hv {
		v, isFloat := val.Value.(float64)
		require.True(t, isFloat)
		require.GreaterOrEqual(t, v, 0.0)
	} else {
		require.FailNow(t, "Did not attach memory metrics")
	}
}
