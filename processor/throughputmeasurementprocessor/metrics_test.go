package throughputmeasurementprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricViews(t *testing.T) {
	expectedViewNames := []string{
		"processor/throughputmeasurement/log_data_size",
		"processor/throughputmeasurement/metric_data_size",
		"processor/throughputmeasurement/trace_data_size",
	}

	views := MetricViews()
	for i, viewName := range expectedViewNames {
		assert.Equal(t, viewName, views[i].Name)
	}
}
