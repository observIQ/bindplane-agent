// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package marshalprocessor

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type marshalProcessor struct {
	logger    *zap.Logger
	marshalTo string
}

func newMarshalProcessor(logger *zap.Logger, cfg *Config) *marshalProcessor {
	return &marshalProcessor{
		logger:    logger,
		marshalTo: cfg.MarshalTo,
	}
}

func (mp *marshalProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resourceLog := ld.ResourceLogs().At(i)
		for j := 0; j < resourceLog.ScopeLogs().Len(); j++ {
			scopeLog := resourceLog.ScopeLogs().At(j)
			for k := 0; k < scopeLog.LogRecords().Len(); k++ {
				logRecord := scopeLog.LogRecords().At(k)
				logBody := logRecord.Body()
				// If body is not a map, skip that log
				if logBody.Type() != pcommon.ValueTypeMap {
					mp.logger.Warn("Non map body not supported", zap.Any("body", logBody))
					continue
				}
				switch strings.ToUpper(mp.marshalTo) {
				case "JSON":
					jsonBody := logBody.AsString()
					logBody.SetStr(jsonBody)
				case "XML":
					return ld, fmt.Errorf("XML not yet supported")
				case "KV":
					kvBody := convertToKV(logBody.Map())
					logBody.SetStr(kvBody)
				default:
					return ld, fmt.Errorf("Unrecognized format to marshal to: %s", mp.marshalTo)
				}
			}
		}
	}

	return ld, nil
}

func convertToKV(logBody pcommon.Map) string {
	var kvStrings []string
	for k, v := range logBody.AsRaw() {
		kvStrings = append(kvStrings, fmt.Sprintf("%s=%v", k, v))
	}

	sort.Strings(kvStrings)

	return strings.Join(kvStrings, " ")
}
