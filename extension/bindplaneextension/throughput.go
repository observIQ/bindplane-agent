package bindplaneextension

import (
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// ThroughputMetricProvider is an interface that marks a component that can be queried for throughput metrics
type ThroughputMetricProvider interface {
	Metrics() ThroughputMetrics
}

type ThroughputMetrics struct {
	logSize     int64
	metricSize  int64
	traceSize   int64
	logCount    int64
	metricCount int64
	traceCount  int64
}

func (tm *ThroughputMetrics) AddLogs(l plog.Logs) {

}

func (tm *ThroughputMetrics) AddMetrics(m pmetric.Metrics) {

}

func (tm *ThroughputMetrics) AddTraces(t ptrace.Traces) {

}
