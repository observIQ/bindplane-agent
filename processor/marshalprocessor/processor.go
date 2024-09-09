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
	logger          *zap.Logger
	marshalTo       string
	kvSeparator     string
	kvPairSeparator string
}

func newMarshalProcessor(logger *zap.Logger, cfg *Config) *marshalProcessor {
	return &marshalProcessor{
		logger:          logger,
		marshalTo:       cfg.MarshalTo,
		kvSeparator:     string(cfg.KVSeparator),
		kvPairSeparator: string(cfg.KVPairSeparator),
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
					kvBody := mp.convertMapToKV(logBody.Map(), false)
					logBody.SetStr(kvBody)
				default:
					return ld, fmt.Errorf("Unrecognized format to marshal to: %s", mp.marshalTo)
				}
			}
		}
	}

	return ld, nil
}

func (mp *marshalProcessor) convertMapToKV(logBody pcommon.Map, inNestedValue bool) string {
	var kvStrings []string

	logBody.Range(func(k string, v pcommon.Value) bool {
		k = mp.escapeAndQuoteKV(k, inNestedValue)

		var vStr string
		switch v.Type() {
		case pcommon.ValueTypeMap:
			vStr = mp.convertMapToKV(v.Map(), true)
			vStr = `[` + vStr + `]`
			if !inNestedValue && (mp.kvSeparator == `=` || mp.kvPairSeparator == `,`) {
				vStr = `"` + vStr + `"`
			}
		default:
			vStr = mp.escapeAndQuoteKV(v.AsString(), inNestedValue)
		}

		if !inNestedValue {
			kvStrings = append(kvStrings, fmt.Sprintf("%s%s%v", k, mp.kvSeparator, vStr))
		} else {
			kvStrings = append(kvStrings, fmt.Sprintf("%s%s%v", k, `=`, vStr))
		}

		return true
	})

	sort.Strings(kvStrings)

	if !inNestedValue {
		return strings.Join(kvStrings, mp.kvPairSeparator)
	} else {
		return strings.Join(kvStrings, `,`)
	}
}

func (mp marshalProcessor) escapeAndQuoteKV(s string, inNestedValue bool) string {
	if !inNestedValue {
		s = strings.ReplaceAll(s, `"`, `\"`)
		if strings.Contains(s, mp.kvPairSeparator) || strings.Contains(s, mp.kvSeparator) {
			s = `"` + s + `"`
		}
	} else {
		s = strings.ReplaceAll(s, `"`, `\\\"`)
		if strings.ContainsAny(s, `=,[]`) {
			s = `\"` + s + `\"`
		}
	}

	return s
}

// func convertMapToString(m map[string]interface{}) string {
// 	var kvPairs []string
// 	var keys []string
// 	for k := range m {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)
