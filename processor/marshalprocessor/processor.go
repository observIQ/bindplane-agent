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
	"encoding/json"
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

// ErrNonMapBodyNotSupported is returned when a log body is not a map
var ErrNonMapBodyNotSupported = fmt.Errorf("Non map body not supported")

func (mp *marshalProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	var errors []error
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resourceLog := ld.ResourceLogs().At(i)
		for j := 0; j < resourceLog.ScopeLogs().Len(); j++ {
			scopeLog := resourceLog.ScopeLogs().At(j)
			for k := 0; k < scopeLog.LogRecords().Len(); k++ {
				logRecord := scopeLog.LogRecords().At(k)
				logBody := logRecord.Body()
				// If body is not a map, skip that log
				if logBody.Type() != pcommon.ValueTypeMap {
					errors = append(errors, ErrNonMapBodyNotSupported)
					continue
				}
				switch strings.ToUpper(mp.marshalTo) {
				case "JSON":
					jsonBody := logBody.AsString()
					logBody.SetStr(jsonBody)
				case "XML":
					return ld, fmt.Errorf("XML not yet supported")
				case "KV":
					jsonBody := logBody.AsString()
					kvBody, err := convertJSONToKV(jsonBody)
					if err != nil {
						errors = append(errors, fmt.Errorf("Error converting to KV: %w", err))
						continue
					}
					logBody.SetStr(kvBody)
				default:
					return ld, fmt.Errorf("Unrecognized format to marshal to: %s", mp.marshalTo)
				}
			}
		}
	}

	if len(errors) == 0 {
		return ld, nil
	}
	var errorStrings []string
	for _, err := range errors {
		errorStrings = append(errorStrings, err.Error())
	}
	return ld, fmt.Errorf("%s", strings.Join(errorStrings, ", "))
}

func convertJSONToKV(jsonBody string) (string, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonBody), &data)
	if err != nil {
		return jsonBody, fmt.Errorf("Error unmarshalling JSON: %w", err)
	}

	var keyValuePairs []string
	for key, value := range data {
		strValue := fmt.Sprintf("%v", value)
		keyValuePairs = append(keyValuePairs, fmt.Sprintf("%s=%s", key, strValue))
	}

	// Ensure consistent order for testing
	sort.Slice(keyValuePairs, func(i, j int) bool {
		return keyValuePairs[i] < keyValuePairs[j]
	})

	kvBody := strings.Join(keyValuePairs, " ")

	return kvBody, nil
}
