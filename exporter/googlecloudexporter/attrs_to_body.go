package googlecloudexporter

import (
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

func logAttrsToBody(ld plog.Logs, keepAttrs map[string]struct{}, retainRawLog bool) {
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		rl := ld.ResourceLogs().At(0)
		for j := 0; j < rl.ScopeLogs().Len(); j++ {
			sl := rl.ScopeLogs().At(j)
			for k := 0; k < sl.LogRecords().Len(); k++ {
				lr := sl.LogRecords().At(k)
				attrsToBody(lr, keepAttrs, retainRawLog)
			}
		}
	}
}

func attrsToBody(lr plog.LogRecord, keepAttrs map[string]struct{}, keepRaw bool) {
	if lr.Body().Type() == pcommon.ValueTypeMap {
		// In this case, the body is already structured.
		// We will ignore the record in this case, and assume it is already in
		// an acceptable format for exporting.
		return
	}

	newBody := pcommon.NewMap()
	newBody.EnsureCapacity(lr.Attributes().Len() + 1)

	lr.Attributes().RemoveIf(func(s string, v pcommon.Value) bool {
		if _, ok := keepAttrs[s]; ok {
			// We should keep this attribute since it's in our keep set
			return false
		}

		if strings.HasPrefix(s, "gcp.") {
			// These are special attributes that are mapped special by the exporter.
			// The exporter ignores these attributes during label mapping, so we will too.
			// https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/blob/de1999d028f2db7630d0cc0306731f6457b6bfae/exporter/collector/logs.go#L415
			return false
		}

		// move to new body
		v.CopyTo(newBody.PutEmpty(s))

		// remove this key from attributes
		return true
	})

	// If the new body would be empty, we'll actually prefer to keep whatever the original body is.
	// This avoids scenarios where the log line gets erased, despite no parsing being applied (e.g. file_logs plugin)
	if newBody.Len() != 0 {
		if keepRaw {
			// retain original log as "raw_log"
			lr.Body().CopyTo(newBody.PutEmpty("raw_log"))
		}
		newBody.CopyTo(lr.Body().SetEmptyMapVal())
	}
}
