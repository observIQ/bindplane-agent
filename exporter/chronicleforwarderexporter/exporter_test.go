// Copyright observIQ, Inc.
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

package chronicleforwarderexporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/plog"
)

func Test_exporter_Capabilities(t *testing.T) {
	exp := &chronicleForwarderExporter{}
	capabilities := exp.Capabilities()
	require.False(t, capabilities.MutatesData)
}

// MockWriter is a mock implementation of io.Writer.
type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func TestLogsDataPusher(t *testing.T) {
	mockWriter := new(MockWriter)
	exporter := &chronicleForwarderExporter{
		writer: mockWriter,
		marshaler: &marshaler{
			cfg: Config{
				ExportType: ExportTypeSyslog,
			},
		},
	}

	mockWriter.On("Write", mock.Anything).Return(0, nil)

	logs := mockLogs([]plog.LogRecord{
		mockLogRecord(t, "Test body", map[string]any{"key1": "value1"}),
	}...)

	err := exporter.logsDataPusher(context.Background(), logs)
	require.NoError(t, err)

	mockWriter.AssertExpectations(t)
}
