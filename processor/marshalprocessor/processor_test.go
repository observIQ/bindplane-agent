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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

func Test_processLogs(t *testing.T) {
	ld := plog.NewLogs()
	ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty().LogRecords().AppendEmpty()

	testCases := []struct {
		desc      string
		serializeTo string
		inputFilePath string
		expected  string
	}{
		{
			desc:      "json",
			serializeTo: "JSON",
			inputFilePath: "json-1.json",
			expected:  "",
		},
		{
			desc:      "kv",
			serializeTo: "KV",
			inputFilePath: "json-1.json",
			expected:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := &Config{
				SerializeTo: tc.serializeTo,
			}

			processor := newSerializeProcessor(zap.NewNop(), cfg)
			actual, err := processor.processLogs(context.Background(), readLogs(t, filepath.Join("testdata", "input", tc.inputFilePath)))
			require.NoError(t, err)
			actualBody := actual.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
			require.Equal(t, tc.expected, actualBody)
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
