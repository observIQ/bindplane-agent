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

package m365receiver // import "github.com/observiq/observiq-otel-collector/receiver/m365receiver"

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/plogtest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestStartPolling(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.TenantID = "testTenantID"
	cfg.Logs.PollInterval = 1 * time.Second

	sink := &consumertest.LogsSink{}
	l := newM365Logs(cfg, receivertest.NewNopCreateSettings(), sink)
	client := &mockLogsClient{}
	l.client = client
	client.On("GetJSON", mock.Anything).Return(client.loadTestLogs(t), nil)
	cancelCtx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel

	err := l.startPolling(cancelCtx)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		return sink.LogRecordCount() > 0
	}, 5*time.Second, 1*time.Second)

	err = l.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestPoll(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.TenantID = "testTenantID"
	cfg.Logs.PollInterval = 1 * time.Second

	sink := &consumertest.LogsSink{}
	rcv := newM365Logs(cfg, receivertest.NewNopCreateSettings(), sink)
	client := &mockLogsClient{}
	rcv.client = client
	client.On("GetJSON", mock.Anything).Return(client.loadTestLogs(t), nil)

	err := rcv.pollLogs(context.Background())
	require.NoError(t, err)

	logs := sink.AllLogs()

	// write logs for each service resource
	// for _, l := range logs {
	// 	audit, exist := l.ResourceLogs().At(0).Resource().Attributes().Get("m365.audit")
	// 	require.True(t, exist)

	// 	marshaler := &plog.JSONMarshaler{}
	// 	lBytes, err := marshaler.MarshalLogs(l)
	// 	require.NoError(t, err)
	// 	err = os.WriteFile(filepath.Join("testdata", "logs", "testPoll", fmt.Sprintf("%s.json", audit.Str())), lBytes, 0666)
	// 	require.NoError(t, err)
	// }

	// compare logs for each service resource
	for _, l := range logs {
		audit, exist := l.ResourceLogs().At(0).Resource().Attributes().Get("m365.audit")
		require.True(t, exist)

		expected, err := ReadLogs(filepath.Join("testdata", "logs", "testPoll", fmt.Sprintf("%s.json", audit.Str())))
		require.NoError(t, err)
		require.NoError(t, plogtest.CompareLogs(expected, l, plogtest.IgnoreObservedTimestamp()))
	}
}

func TestTransformLogs(t *testing.T) {
	sink := &consumertest.LogsSink{}
	now := pcommon.NewTimestampFromTime(time.Time{})
	audit := auditMetaData{"General", "Audit.General", true}
	file := filepath.Join("testdata", "logs", "transform-test.json")
	logData := []jsonLogs{
		{
			Workload:     "testWorkload",
			UserId:       "testUserID",
			UserType:     0,
			CreationTime: "2023-05-10T09:07:33",
			Id:           "testID",
			Operation:    "testOperation",
			ResultStatus: "testResultStatus",
		},
		{
			Workload:     "testWorkload2",
			UserId:       "testUserID2",
			UserType:     0,
			CreationTime: "2023-05-10T09:07:33",
			Id:           "testID2",
			Operation:    "testOperation2",
			ResultStatus: "testResultStatus2",
		},
	}

	cfg := NewFactory().CreateDefaultConfig().(*Config)
	l := newM365Logs(cfg, receivertest.NewNopCreateSettings(), sink)
	result := l.transformLogs(now, &audit, logData)

	// write logs to file
	// marshaler := &plog.JSONMarshaler{}
	// lBytes, err := marshaler.MarshalLogs(result)
	// require.NoError(t, err)
	// err = os.WriteFile(file, lBytes, 0666)
	// require.NoError(t, err)

	// compare logs
	expected, err := ReadLogs(file)
	require.NoError(t, err)
	require.NoError(t, plogtest.CompareLogs(expected, result, plogtest.IgnoreObservedTimestamp()))
}

func TestParseOptionalAttributes(t *testing.T) {
	m := pcommon.NewMap()
	log := jsonLogs{
		Workload:     "testWorkload",
		ResultStatus: "",
	}
	parseOptionalAttributes(&m, &log)
	w, e := m.Get("workload")
	require.True(t, e)
	require.Equal(t, "testWorkload", w.AsString())
	w, e = m.Get("result_status")
	require.False(t, e)
}

func ReadLogs(filepath string) (plog.Logs, error) {
	expectedFileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return plog.Logs{}, nil
	}
	unmarshaler := &plog.JSONUnmarshaler{}
	return unmarshaler.UnmarshalLogs(expectedFileBytes)
}

type mockLogsClient struct {
	mock.Mock
}

func (mc *mockLogsClient) loadTestLogs(t *testing.T) []jsonLogs {
	testLogs := filepath.Join("testdata", "logs", "testPoll", "input.json")
	logBytes, err := os.ReadFile(testLogs)
	require.NoError(t, err)

	var logs []jsonLogs
	err = json.Unmarshal(logBytes, &logs)
	require.NoError(t, err)
	return logs
}

func (mc *mockLogsClient) GetJSON(endpoint string) ([]jsonLogs, error) {
	args := mc.Called(endpoint)
	return args.Get(0).([]jsonLogs), args.Error(1)
}

func (mc *mockLogsClient) GetToken() error {
	return nil
}

func (mc *mockLogsClient) shutdown() error {
	return nil
}
