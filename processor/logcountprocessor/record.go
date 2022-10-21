package logcountprocessor

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// Record is an observed log record.
type Record struct {
	Resource   map[string]interface{}
	Attributes map[string]interface{}
}

// NewRecord creates a new log record.
func NewRecord(resource, attributes map[string]interface{}) Record {
	return Record{
		resource,
		attributes,
	}
}

// ResourceAsMap returns the log record's resource as a pcommon.Map
func (l Record) ResourceAsMap() pcommon.Map {
	m := pcommon.NewMap()
	m.FromRaw(l.Resource)
	return m
}

// AttributesAsMap returns the log record's attributes as a pcommon.Map
func (l Record) AttributesAsMap() pcommon.Map {
	m := pcommon.NewMap()
	m.FromRaw(l.Attributes)
	return m
}
