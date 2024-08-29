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

package serializeprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type serializeProcessor struct {
	logger          *zap.Logger
	serializeTo	string
}

func newSerializeProcessor(logger *zap.Logger, cfg *Config) *serializeProcessor {
	return &serializeProcessor{
		logger:          logger,
		serializeTo: cfg.SerializeTo,
	}
}

// The use of this processor assumes that the body of the log record has already been parsed. 
func (sp *serializeProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	var errors []error
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		for j := 0; j < ld.ResourceLogs().At(i).ScopeLogs().Len(); j++ {
			for k := 0; k < ld.ResourceLogs().At(i).ScopeLogs().At(j).LogRecords().Len(); k++ {
				var body = ld.ResourceLogs().At(i).ScopeLogs().At(j).LogRecords().At(k).Body()
				switch strings.ToUpper(sp.serializeTo) {
				case "JSON":
					var jsonBody = body.AsString()
					ld.ResourceLogs().At(i).ScopeLogs().At(j).LogRecords().At(k).Body().SetStr(jsonBody)
				// case "XML": 
				case "KV":
					var jsonBody = body.AsString()
					kvBody, err := convertJSONToKV(jsonBody)
					if err != nil {
						errors = append(errors, err)
						continue
					}
					ld.ResourceLogs().At(i).ScopeLogs().At(j).LogRecords().At(k).Body().SetStr(kvBody)
				default:
					errors = append(errors, fmt.Errorf("No recognized serialization option provided"))
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
	return ld, fmt.Errorf(strings.Join(errorStrings, ", "))
}

// Assumes data is flattened for nice format
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

	kvBody := strings.Join(keyValuePairs, " ")

	return kvBody, nil
}