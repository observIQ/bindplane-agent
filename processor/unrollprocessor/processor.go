package unrollprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

type unrollProcessor struct {
	cfg *Config
}

// newUnrollProcessor returns a new unrollProcessor.
func newUnrollProcessor(config *Config) (*unrollProcessor, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &unrollProcessor{
		cfg: config,
	}, nil
}

// ProcessLogs implements the processor interface
func (p *unrollProcessor) ProcessLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	newLogs := plog.NewLogs()
	resourceLogs := ld.ResourceLogs()
	var errs error
	for i := 0; i < resourceLogs.Len(); i++ {
		resourceLog := resourceLogs.At(i)
		sl := resourceLog.ScopeLogs()
		p.siftScopeLogs(sl)
	}
	resourceLogs.CopyTo(newLogs.ResourceLogs())
	return newLogs, errs
}

// siftScopeLogs will iterate through scope logs and copy over new logs from siftLogRecords
func (p *unrollProcessor) siftScopeLogs(sl plog.ScopeLogsSlice) {
	for i := 0; i < sl.Len(); i++ {
		scopeLogs := sl.At(i)
		l := scopeLogs.LogRecords()
		p.siftLogRecords(l).CopyTo(scopeLogs.LogRecords())
	}
}

// siftLogRecords goes through each log record in the slice and will add new logs if there is a slice body
// returns a log record slice with any logs that were not a slice body and new logs from the slice body
func (p *unrollProcessor) siftLogRecords(logs plog.LogRecordSlice) plog.LogRecordSlice {
	newLogRecords := plog.NewLogRecordSlice()
	for i := 0; i < logs.Len(); i++ {
		logRecord := logs.At(i)
		// if the body is not a slice, retain the original log
		if logRecord.Body().Type() != pcommon.ValueTypeSlice {
			logRecord.CopyTo(newLogRecords.AppendEmpty())
			continue
		}

		// get length of slice, make that many new log records
		// TODO: optimized: n-1 copies, modify original log record
		for j := 0; j < logRecord.Body().Slice().Len(); j++ {
			expansion := logRecord.Body().Slice().At(j)
			newLogRecord := plog.NewLogRecord()
			logRecord.CopyTo(newLogRecord)
			if p.cfg.UnrollKey == "" {
				p.setBody(newLogRecord, expansion)
			} else {
				p.addToNewMap(newLogRecord, expansion)
			}
			newLogRecord.CopyTo(newLogRecords.AppendEmpty())
		}

	}
	return newLogRecords
}

// addToNewMap will add a new map body to the log record with the unroll key and the provided value
func (p *unrollProcessor) addToNewMap(newLogRecord plog.LogRecord, value pcommon.Value) {
	m := newLogRecord.Body().SetEmptyMap()
	switch value.Type() {
	case pcommon.ValueTypeStr:
		m.PutStr(p.cfg.UnrollKey, value.Str())
	case pcommon.ValueTypeInt:
		m.PutInt(p.cfg.UnrollKey, value.Int())
	case pcommon.ValueTypeDouble:
		m.PutDouble(p.cfg.UnrollKey, value.Double())
	case pcommon.ValueTypeBool:
		m.PutBool(p.cfg.UnrollKey, value.Bool())
	case pcommon.ValueTypeMap:
		value.Map().CopyTo(m.PutEmptyMap(p.cfg.UnrollKey))
	case pcommon.ValueTypeSlice:
		value.Slice().CopyTo(m.PutEmptySlice(p.cfg.UnrollKey))
	}
}

// setBody will set the body of the log record to the provided value
func (p *unrollProcessor) setBody(newLogRecord plog.LogRecord, expansion pcommon.Value) {
	switch expansion.Type() {
	case pcommon.ValueTypeStr:
		newLogRecord.Body().SetStr(expansion.Str())
	case pcommon.ValueTypeInt:
		newLogRecord.Body().SetInt(expansion.Int())
	case pcommon.ValueTypeDouble:
		newLogRecord.Body().SetDouble(expansion.Double())
	case pcommon.ValueTypeBool:
		newLogRecord.Body().SetBool(expansion.Bool())
	case pcommon.ValueTypeMap:
		expansion.Map().CopyTo(newLogRecord.Body().SetEmptyMap())
	case pcommon.ValueTypeSlice:
		expansion.Slice().CopyTo(newLogRecord.Body().SetEmptySlice())
	}
}
