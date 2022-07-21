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
	Record   interface{} `json:"record"`
	Sessions []string    `json:"sessions"`
}

// NewMetricsMessage creates a new message containing metrics
func NewMetricsMessage(record MetricRecord, sessions []string) Message {
	return Message{
		Type:     MetricsType,
		Sessions: sessions,
		Record:   record,
	}
}

// NewLogsMessage creates a new message containing logs
func NewLogsMessage(record LogRecord, sessions []string) Message {
	return Message{
		Type:     LogsType,
		Sessions: sessions,
		Record:   record,
	}
}

// NewTracesMessage creates a new message containing traces
func NewTracesMessage(record TraceRecord, sessions []string) Message {
	return Message{
		Type:     TracesType,
		Sessions: sessions,
		Record:   record,
	}
}
