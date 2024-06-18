//go:build !bindplane

package collector

// ResetMeasurements resets the registered throughput measurements
// It is a no-op for non-bindplane agents
func ResetMeasurements() {}
