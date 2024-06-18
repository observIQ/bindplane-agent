//go:build !bindplane

package throughputmeasurementprocessor

import (
	"fmt"

	"github.com/observiq/bindplane-agent/internal/measurements"
	"go.opentelemetry.io/collector/component"
)

// GetThroughputRegistry returns the throughput registry that should be registered to based on the component ID.
// nil, nil may be returned by this function. In this case, the processor should not register it's throughput measurements anywhere.
func GetThroughputRegistry(host component.Host, bindplane component.ID) (measurements.ThroughputMeasurementsRegistry, error) {
	var emptyComponentID component.ID
	if bindplane == emptyComponentID {
		// No bindplane component referenced, so we won't register our measurements anywhere.
		return nil, nil
	}

	ext, ok := host.GetExtensions()[bindplane]
	if !ok {
		return nil, fmt.Errorf("bindplane extension %q does not exist", bindplane)
	}

	registry, ok := ext.(measurements.ThroughputMeasurementsRegistry)
	if !ok {
		return nil, fmt.Errorf("extension %q is not an throughput message registry", bindplane)
	}

	return registry, nil
}
