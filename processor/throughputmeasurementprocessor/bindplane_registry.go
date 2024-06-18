//go:build bindplane

package throughputmeasurementprocessor

import (
	"github.com/observiq/bindplane-agent/internal/measurements"
	"go.opentelemetry.io/collector/component"
)

// GetThroughputRegistry returns the throughput registry that should be registered to based on the component ID.
// nil, nil may be returned by this function. In this case, the processor should not register it's throughput measurements anywhere.
func GetThroughputRegistry(host component.Host, bindplane component.ID) (measurements.ThroughputMeasurementsRegistry, error) {
	return measurements.BindplaneAgentThroughputMeasurementsRegistry, nil
}
