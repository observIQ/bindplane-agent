// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chronicleexporter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/observiq/bindplane-agent/expr"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

const logTypeField = `attributes["log_type"]`

var supportedLogTypes = map[string]string{
	"windows_event.security":    "WINEVTLOG",
	"windows_event.custom":      "WINEVTLOG",
	"windows_event.application": "WINEVTLOG",
	"windows_event.system":      "WINEVTLOG",
	"sql_server":                "MICROSOFT_SQL",
}

type marshaler struct {
	cfg          Config
	teleSettings component.TelemetrySettings
}

type payload struct {
	Entries    []entry `json:"entries"`
	CustomerID string  `json:"customer_id"`
	LogType    string  `json:"log_type"`
	Namespace  string  `json:"namespace"`
}

type entry struct {
	LogText   string `json:"log_text"`
	Timestamp string `json:"ts_rfc3339"`
}

type logMarshaler interface {
	MarshalRawLogs(ctx context.Context, ld plog.Logs) ([]payload, error)
}

func newMarshaler(cfg Config, teleSettings component.TelemetrySettings) *marshaler {
	return &marshaler{
		cfg:          cfg,
		teleSettings: teleSettings,
	}
}

func (m *marshaler) MarshalRawLogs(ctx context.Context, ld plog.Logs) ([]payload, error) {
	rawLogs, err := m.extractRawLogs(ctx, ld)
	if err != nil {
		return nil, fmt.Errorf("extract raw logs: %w", err)
	}

	return m.constructPayloads(rawLogs), nil
}

func (m *marshaler) extractRawLogs(ctx context.Context, ld plog.Logs) (map[string][]entry, error) {
	entries := make(map[string][]entry)

	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resourceLog := ld.ResourceLogs().At(i)
		for j := 0; j < resourceLog.ScopeLogs().Len(); j++ {
			scopeLog := resourceLog.ScopeLogs().At(j)
			for k := 0; k < scopeLog.LogRecords().Len(); k++ {
				logRecord := scopeLog.LogRecords().At(k)
				rawLog, logType, err := m.processLogRecord(ctx, logRecord, scopeLog.Scope(), resourceLog.Resource())
				if err != nil {
					m.teleSettings.Logger.Error("Error processing log record", zap.Error(err))
					continue
				}

				if rawLog == "" {
					continue
				}

				var timestamp time.Time
				if logRecord.Timestamp() != 0 {
					timestamp = logRecord.Timestamp().AsTime()
				} else {
					timestamp = logRecord.ObservedTimestamp().AsTime()
				}

				entries[logType] = append(entries[logType], entry{
					LogText:   rawLog,
					Timestamp: timestamp.Format(time.RFC3339Nano),
				})
			}
		}
	}

	return entries, nil
}

func (m *marshaler) processLogRecord(ctx context.Context, logRecord plog.LogRecord, scope pcommon.InstrumentationScope, resource pcommon.Resource) (string, string, error) {
	rawLog, err := m.getRawLog(ctx, logRecord, scope, resource)
	if err != nil {
		return "", "", err
	}

	logType, err := m.getLogType(ctx, logRecord, scope, resource)
	if err != nil {
		return "", "", err
	}

	return rawLog, logType, nil
}

func (m *marshaler) getRawLog(ctx context.Context, logRecord plog.LogRecord, scope pcommon.InstrumentationScope, resource pcommon.Resource) (string, error) {
	if m.cfg.RawLogField == "" {
		entireLogRecord := map[string]any{
			"body":                logRecord.Body().Str(),
			"attributes":          logRecord.Attributes().AsRaw(),
			"resource_attributes": resource.Attributes().AsRaw(),
		}

		bytesLogRecord, err := json.Marshal(entireLogRecord)
		if err != nil {
			return "", fmt.Errorf("marshal log record: %w", err)
		}

		return string(bytesLogRecord), nil
	}
	return m.getRawField(ctx, m.cfg.RawLogField, logRecord, scope, resource)
}

func (m *marshaler) getLogType(ctx context.Context, logRecord plog.LogRecord, scope pcommon.InstrumentationScope, resource pcommon.Resource) (string, error) {
	if m.cfg.OverrideLogType {
		logType, err := m.getRawField(ctx, logTypeField, logRecord, scope, resource)
		if err != nil || logType == "" {
			return m.cfg.LogType, err
		}

		if chronicleLogType, ok := supportedLogTypes[logType]; ok {
			return chronicleLogType, nil
		}
	}

	return m.cfg.LogType, nil
}

func (m *marshaler) getRawField(ctx context.Context, field string, logRecord plog.LogRecord, scope pcommon.InstrumentationScope, resource pcommon.Resource) (string, error) {
	lrExpr, err := expr.NewOTTLLogRecordExpression(field, m.teleSettings)
	if err != nil {
		return "", fmt.Errorf("raw_log_field is invalid: %s", err)
	}
	tCtx := ottllog.NewTransformContext(logRecord, scope, resource)

	lrExprResult, err := lrExpr.Execute(ctx, tCtx)
	if err != nil {
		return "", fmt.Errorf("execute log record expression: %w", err)
	}

	if lrExprResult == nil {
		return "", nil
	}

	switch result := lrExprResult.(type) {
	case string:
		return result, nil
	case pcommon.Map:
		bytes, err := json.Marshal(result.AsRaw())
		if err != nil {
			return "", fmt.Errorf("marshal log record expression result: %w", err)
		}
		return string(bytes), nil
	default:
		return "", fmt.Errorf("unsupported log record expression result type: %T", lrExprResult)
	}
}

func (m *marshaler) constructPayloads(rawLogs map[string][]entry) []payload {
	payloads := make([]payload, 0, len(rawLogs))
	for logType, entries := range rawLogs {
		if len(entries) > 0 {
			p := payload{
				Entries:    entries,
				CustomerID: m.cfg.CustomerID,
				LogType:    logType,
			}

			if m.cfg.Namespace != "" {
				p.Namespace = m.cfg.Namespace
			}

			payloads = append(payloads, p)
		}
	}
	return payloads
}
