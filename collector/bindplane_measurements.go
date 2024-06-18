// go:build bindplane
package collector

import "github.com/observiq/bindplane-agent/internal/measurements"

// ResetMeasurements resets the registered throughput measurements
func ResetMeasurements() {
	measurements.BindplaneAgentThroughputMeasurementsRegistry.Reset()
}
