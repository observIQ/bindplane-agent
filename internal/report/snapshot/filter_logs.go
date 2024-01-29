package snapshot

import (
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

// filterLogs filters the logs by the given query and timestamp.
// The returned payload cannot be assumed to be a copy, so it should not be modified.
func filterLogs(logs plog.Logs, query *string, minTimestamp *time.Time) plog.Logs {
	// No filters specified, filtered logs are trivially the same as input logs
	if query == nil && minTimestamp == nil {
		return logs
	}

	filteredLogs := plog.NewLogs()

	resourceLogs := logs.ResourceLogs()
	for i := 0; i < resourceLogs.Len(); i++ {
		filteredResourceLogs := filterResourceLogs(resourceLogs.At(i), query, minTimestamp)

		// Don't append empty resource logs
		if filteredResourceLogs.ScopeLogs().Len() != 0 {
			filteredResourceLogs.CopyTo(filteredLogs.ResourceLogs().AppendEmpty())
		}
	}

	return filteredLogs
}

func filterResourceLogs(resourceLog plog.ResourceLogs, query *string, minTimestamp *time.Time) plog.ResourceLogs {
	filteredResourceLogs := plog.NewResourceLogs()

	// Copy old resource to filtered resource
	resource := resourceLog.Resource()
	resource.CopyTo(filteredResourceLogs.Resource())

	// Apply query to resource
	queryMatchesResource := true // default to true if no query specified
	if query != nil {
		queryMatchesResource = queryMatchesMap(resource.Attributes(), *query)
	}

	scopeLogs := resourceLog.ScopeLogs()
	for i := 0; i < scopeLogs.Len(); i++ {
		filteredScopeLogs := filterScopeLogs(resourceLog.ScopeLogs().At(i), queryMatchesResource, query, minTimestamp)

		// Don't append empty scope logs
		if filteredScopeLogs.LogRecords().Len() != 0 {
			filteredScopeLogs.CopyTo(filteredResourceLogs.ScopeLogs().AppendEmpty())
		}
	}

	return filteredResourceLogs
}

// filterScopeLogs filters out logs that do not match the query and minTimestamp, returning a new plog.ScopeLogs without the filtered records.
// queryMatchesResource indicates if the query string matches the resource associated with this ScopeLogs.
func filterScopeLogs(scopeLogs plog.ScopeLogs, queryMatchesResource bool, query *string, minTimestamp *time.Time) plog.ScopeLogs {
	filteredLogRecords := plog.NewScopeLogs()
	logRecords := scopeLogs.LogRecords()
	for i := 0; i < logRecords.Len(); i++ {
		log := logRecords.At(i)
		if logMatches(log, queryMatchesResource, query, minTimestamp) {
			log.CopyTo(filteredLogRecords.LogRecords().AppendEmpty())
		}
	}

	return filteredLogRecords
}

// logMatches returns true if the query matches either the resource or log record, AND the min timestamp.
func logMatches(l plog.LogRecord, queryMatchesResource bool, query *string, minTimestamp *time.Time) bool {
	queryMatchesLog := true // default to true if no query specified
	// Skip this check if we already know the query matches the resource
	if !queryMatchesResource && query != nil {
		queryMatchesLog = logMatchesQuery(l, *query)
	}

	timestampMatches := true // default to true if no timestamp specified
	if minTimestamp != nil {
		timestampMatches = logMatchesTimestamp(l, *minTimestamp)
	}

	queryMatches := queryMatchesResource || queryMatchesLog

	return queryMatches && timestampMatches
}

// logMatchesTimestamp determines if the log came after the provided timestamp
func logMatchesTimestamp(l plog.LogRecord, minTime time.Time) bool {
	return l.ObservedTimestamp() >= pcommon.NewTimestampFromTime(minTime)
}

// logMatchesQuery determines if the given log record matches the given query string
func logMatchesQuery(l plog.LogRecord, query string) bool {
	if queryMatchesMap(l.Attributes(), query) {
		return true
	}

	if queryMatchesValue(l.Body(), query) {
		return true
	}

	return false
}
