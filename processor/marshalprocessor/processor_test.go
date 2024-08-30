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
	"path/filepath"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/plogtest"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_processLogs(t *testing.T) {
	testCases := []struct {
		desc                   string
		marshalTo              string
		inputFilePath          string
		expectedOutputFilePath string
	}{
		{
			desc:                   "Valid - Parsed body to JSON",
			marshalTo:              "JSON",
			inputFilePath:          "parsed-log-input.json",
			expectedOutputFilePath: `parsed-log-json-output.json`,
		},
		{
			desc:                   "Invalid - String body to JSON",
			marshalTo:              "JSON",
			inputFilePath:          "string-log-input.json",
			expectedOutputFilePath: "string-log-output.json",
		},
		{
			desc:                   "Invalid - String body to KV",
			marshalTo:              "KV",
			inputFilePath:          "string-log-input.json",
			expectedOutputFilePath: "string-log-output.json",
		},
		{
			desc:                   "Invalid - String body to XML",
			marshalTo:              "XML",
			inputFilePath:          "string-log-input.json",
			expectedOutputFilePath: "string-log-output.json",
		},
		{
			desc:                   "Valid - Parsed and flattened body to KV",
			marshalTo:              "KV",
			inputFilePath:          "parsed-flattened-log-input.json",
			expectedOutputFilePath: "parsed-flattened-log-kv-output.json",
		},
		{
			desc:                   "Valid - Parsed nested body to KV", // not recommended to use this unflattened format but technically valid
			marshalTo:              "KV",
			inputFilePath:          "parsed-log-input.json",
			expectedOutputFilePath: "parsed-log-kv-output.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := &Config{
				MarshalTo: tc.marshalTo,
			}

			processor := newMarshalProcessor(zap.NewNop(), cfg)
			inputlogs, err := golden.ReadLogs(filepath.Join("testdata", "input", tc.inputFilePath))
			require.NoError(t, err)
			actual, err := processor.processLogs(context.Background(), inputlogs)
			require.NoError(t, err)
			expectedOutput, err := golden.ReadLogs(filepath.Join("testdata", "output", tc.expectedOutputFilePath))
			require.NoError(t, err)

			require.NoError(t, plogtest.CompareLogs(expectedOutput, actual))
		})
	}
}
