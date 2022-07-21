package bindplaneexporter

import "go.opentelemetry.io/collector/consumer"

// Exporter is the bindplane exporter
type Exporter struct {
	consumeMetrics consumer.ConsumeMetricsFunc
	consumeLogs    consumer.ConsumeLogsFunc
	consumeTraces  consumer.ConsumeTracesFunc
}
