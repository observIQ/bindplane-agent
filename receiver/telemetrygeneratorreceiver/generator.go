package telemetrygeneratorreceiver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type generator interface {
	// Passes in a component.Type which has a value of logs, metrics, traces. Returns true or false if this generator supports it.
	SupportsType(component.Type) bool
	GenerateMetrics() pmetric.Metrics
	GenerateLogs() plog.Logs
	GenerateTraces() ptrace.Traces
}
