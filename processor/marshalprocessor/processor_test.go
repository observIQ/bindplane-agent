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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

func Test_processLogs(t *testing.T) {
	testCases := []struct {
		desc      string
		marshalTo string
		inputFilePath string
		expected  string
		expectedErr error
	}{
		{
			desc:      "Valid - Parsed body to JSON",
			marshalTo: "JSON",
			inputFilePath: "parsed-log-1.json",
			expected:  `{"bindplane-otel-attributes":{"baba":"you","host":"myhost"},"name":"test","nested":{"n1":1,"n2":2},"severity":155}`,
			expectedErr: nil,
		},
		{
			desc:      "Invalid - String body to JSON",
			marshalTo: "JSON",
			inputFilePath: "string-log-1.json",
			expected:  "",
			expectedErr: ErrStringBodyNotSupported,
		},
		{
			desc:      "Invalid - String body to KV",
			marshalTo: "KV",
			inputFilePath: "string-log-1.json",
			expected:  "",
			expectedErr: ErrStringBodyNotSupported,
		},
		{
			desc:      "Invalid - String body to XML",
			marshalTo: "XML",
			inputFilePath: "string-log-1.json",
			expected:  "",
			expectedErr: ErrStringBodyNotSupported,
		},
		{
			desc:      "Valid - Parsed and flattened body to KV",
			marshalTo: "KV",
			inputFilePath: "parsed-flattened-log-1.json",
			expected:  "bindplane-otel-attributes-baba=you bindplane-otel-attributes-host=myhost name=test nested-n1=1 nested-n2=2 severity=155",
			expectedErr: nil,
		},
		{
			desc:      "Valid - Parsed nested body to KV", // not recommended to use this unflattened format but technically valid
			marshalTo: "KV",
			inputFilePath: "parsed-log-1.json",
			expected:  "bindplane-otel-attributes=map[baba:you host:myhost] name=test nested=map[n1:1 n2:2] severity=155",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := &Config{
				MarshalTo: tc.marshalTo,
			}

			processor := newMarshalProcessor(zap.NewNop(), cfg)
			actual, err := processor.processLogs(context.Background(), readLogs(t, filepath.Join("testdata", "input", tc.inputFilePath)))
			require.Equal(t, tc.expectedErr, err)
			if err == nil {
				actualBody := actual.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
				require.Equal(t, tc.expected, actualBody)
			}
		})
	}
}

func readLogs(t *testing.T, path string) plog.Logs {
	t.Helper()

	b, err := os.ReadFile(path)
	require.NoError(t, err)

	unmarshaller := plog.JSONUnmarshaler{}
	l, err := unmarshaller.UnmarshalLogs(b)
	require.NoError(t, err)

	return l
}
