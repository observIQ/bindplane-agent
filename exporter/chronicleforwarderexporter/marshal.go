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

package chronicleforwarderexporter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/observiq/bindplane-otel-collector/internal/expr"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

const logTypeField = `attributes["log_type"]`

type marshaler struct {
	cfg          Config
	teleSettings component.TelemetrySettings
}

type logMarshaler interface {
	MarshalRawLogs(ctx context.Context, ld plog.Logs) ([]string, error)
}

func newMarshaler(cfg Config, teleSettings component.TelemetrySettings) *marshaler {
	return &marshaler{
		cfg:          cfg,
		teleSettings: teleSettings,
	}
}

func (m *marshaler) MarshalRawLogs(ctx context.Context, ld plog.Logs) ([]string, error) {
	var rawLogs []string

	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resourceLog := ld.ResourceLogs().At(i)
		for j := 0; j < resourceLog.ScopeLogs().Len(); j++ {
			scopeLog := resourceLog.ScopeLogs().At(j)
			for k := 0; k < scopeLog.LogRecords().Len(); k++ {
				logRecord := scopeLog.LogRecords().At(k)

				// Escape any unescaped newlines (if body is a string)
				logBody := logRecord.Body()
				if logBody.Type() == pcommon.ValueTypeStr {
					logBody.SetStr(strings.ReplaceAll(logBody.AsString(), "\n", "\\n"))
				}

				rawLog, err := m.getRawLog(ctx, logRecord, scopeLog, resourceLog)
				if err != nil {
					m.teleSettings.Logger.Error("Error processing log record", zap.Error(err))
					continue
				}

				if rawLog == "" {
					m.teleSettings.Logger.Error("Error processing log record: raw log is empty")
					continue
				}

				rawLogs = append(rawLogs, rawLog)
			}
		}
	}

	return rawLogs, nil
}

func (m *marshaler) getRawLog(ctx context.Context, logRecord plog.LogRecord, scope plog.ScopeLogs, resource plog.ResourceLogs) (string, error) {
	if m.cfg.RawLogField == "" {
		entireLogRecord := map[string]any{
			"body":                logRecord.Body().AsRaw(),
			"attributes":          logRecord.Attributes().AsRaw(),
			"resource_attributes": resource.Resource().Attributes().AsRaw(),
		}

		bytesLogRecord, err := json.Marshal(entireLogRecord)
		if err != nil {
			return "", fmt.Errorf("marshal log record: %w", err)
		}

		return string(bytesLogRecord), nil
	}
	return m.getRawField(ctx, m.cfg.RawLogField, logRecord, scope, resource)
}

func (m *marshaler) getRawField(ctx context.Context, field string, logRecord plog.LogRecord, scope plog.ScopeLogs, resource plog.ResourceLogs) (string, error) {
	lrExpr, err := expr.NewOTTLLogRecordExpression(field, m.teleSettings)
	if err != nil {
		return "", fmt.Errorf("raw_log_field is invalid: %s", err)
	}
	tCtx := ottllog.NewTransformContext(logRecord, scope.Scope(), resource.Resource(), scope, resource)

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
