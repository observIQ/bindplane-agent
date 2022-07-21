package bindplaneexporter

const (
	// MetricsType indicates that a message contains metrics
	MetricsType = "metrics"

	// LogsType indicates that a message contains logs
	LogsType = "logs"

	// TracesType indicates that a message contains traces
	TracesType = "traces"
)

// Message is a message sent to bindplane from the exporter
type Message struct {
	Type     string        `json:"type"`
	Sessions []string      `json:"sessions"`
	Records  []interface{} `json:"records"`
}
