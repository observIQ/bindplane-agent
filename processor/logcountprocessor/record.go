package logcountprocessor

import "go.opentelemetry.io/collector/pdata/plog"

const (
	// bodyField is the name of the field containing the log body.
	bodyField = "body"

	// resourceField is the name of the field containing the resource attributes.
	resourceField = "resource"

	// attributesField is the name of the field containing the log attributes.
	attributesField = "attributes"

	// severityField is the name of the field containing the log severity.
	severityField = "severity"
)

// Record is the simplified representation of a log record.
type Record = map[string]interface{}

// convertToRecords converts plog.Logs to a slice of records.
func convertToRecords(logs plog.Logs) []Record {
	records := make([]Record, 0, logs.ResourceLogs().Len())

	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		resourceLogs := logs.ResourceLogs().At(i)
		resource := resourceLogs.Resource().Attributes().AsRaw()
		for j := 0; j < resourceLogs.ScopeLogs().Len(); j++ {
			logs := resourceLogs.ScopeLogs().At(j).LogRecords()
			for k := 0; k < logs.Len(); k++ {
				log := logs.At(k)
				records = append(records, convertToRecord(log, resource))
			}
		}
	}

	return records
}

// convertToRecord converts a log record to a simplified representation.
func convertToRecord(log plog.LogRecord, resource map[string]interface{}) Record {
	return Record{
		resourceField:   resource,
		attributesField: log.Attributes().AsRaw(),
		bodyField:       log.Body().AsRaw(),
		severityField:   log.SeverityText(),
	}
}
