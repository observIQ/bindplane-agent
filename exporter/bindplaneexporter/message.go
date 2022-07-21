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
	Type     string      `json:"type"`
	Records  interface{} `json:"records"`
	Sessions []string    `json:"sessions"`
}

// NewMetricsMessage creates a new message containing metrics
func NewMetricsMessage(records []MetricRecord, sessions []string) Message {
	return Message{
		Type:     MetricsType,
		Sessions: sessions,
		Records:  records,
	}
}

// NewLogsMessage creates a new message containing logs
func NewLogsMessage(records []LogRecord, sessions []string) Message {
	return Message{
		Type:     LogsType,
		Sessions: sessions,
		Records:  records,
	}
}

// NewTracesMessage creates a new message containing traces
func NewTracesMessage(records []TraceRecord, sessions []string) Message {
	return Message{
		Type:     TracesType,
		Sessions: sessions,
		Records:  records,
	}
}
