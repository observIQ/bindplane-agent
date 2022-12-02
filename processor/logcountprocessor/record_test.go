package logcountprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/plog"
)

func TestConvertToRecords(t *testing.T) {
	testResource := map[string]interface{}{
		"resource": "attributes",
	}
	testAttrs := map[string]interface{}{
		"attributes": "attributes",
	}
	testBody := "body"
	testSeverity := "info"

	logs := plog.NewLogs()
	resourceLogs := logs.ResourceLogs().AppendEmpty()
	resourceLogs.Resource().Attributes().FromRaw(testResource)

	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
	logRecords := scopeLogs.LogRecords()
	log1 := logRecords.AppendEmpty()
	log1.Body().SetStr(testBody)
	log1.SetSeverityText(testSeverity)
	log1.Attributes().FromRaw(testAttrs)

	records := convertToRecords(logs)
	require.Len(t, records, 1)
	require.Equal(t, map[string]interface{}{
		resourceField:   testResource,
		attributesField: testAttrs,
		bodyField:       testBody,
		severityField:   testSeverity,
	}, records[0])
}
