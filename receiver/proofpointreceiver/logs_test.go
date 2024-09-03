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

package proofpointreceiver

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/plogtest"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.uber.org/zap"
)

type mockHTTPClient struct {
	mockDo func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.mockDo(req)
}

func (m *mockHTTPClient) CloseIdleConnections() {}

func TestStartShutdown(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	recv, err := newProofpointLogsReceiver(cfg, zap.NewNop(), consumertest.NewNop())
	require.NoError(t, err)

	err = recv.Start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestPollBasic(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Principal = "mockPrincipal"
	cfg.Secret = "mockSecret"

	sink := &consumertest.LogsSink{}
	recv, err := newProofpointLogsReceiver(cfg, zap.NewNop(), sink)
	require.NoError(t, err)

	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			require.Contains(t, req.URL.String(), "sinceTime=")
			require.NotContains(t, req.URL.String(), "sinceTime=1888-06-24")
			require.Contains(t, req.URL.String(), "format=json")
			return mockAPIResponseOK(t), nil
		},
	}

	err = recv.poll(context.Background())
	require.NoError(t, err)

	logs := sink.AllLogs()
	log := logs[0]

	// golden.WriteLogs(t, "testdata/plog.yaml", log)

	expected, err := golden.ReadLogs("testdata/plog.yaml")
	require.NoError(t, err)
	require.NoError(t, plogtest.CompareLogs(expected, log, plogtest.IgnoreObservedTimestamp()))

	require.Equal(t, 4, log.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().Len())

	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			require.Contains(t, req.URL.String(), "sinceTime=1888-06-24T21%3A36%3A00Z")
			require.Contains(t, req.URL.String(), "format=json")
			return mockAPIResponseOK(t), nil
		},
	}
	err = recv.poll(context.Background())
	require.NoError(t, err)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestPoll429TooManyRequests(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Principal = "mockPrincipal"
	cfg.Secret = "mockSecret"

	sink := &consumertest.LogsSink{}
	recv, err := newProofpointLogsReceiver(cfg, zap.NewNop(), sink)
	require.NoError(t, err)

	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			require.Contains(t, req.URL.String(), "sinceTime=")
			require.NotContains(t, req.URL.String(), "sinceTime=1888-06-24")
			require.Contains(t, req.URL.String(), "format=json")
			return mockAPIResponse429(), nil
		},
	}

	err = recv.poll(context.Background())
	require.Error(t, err)

	logs := sink.AllLogs()
	require.Empty(t, logs)

	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			require.Contains(t, req.URL.String(), "sinceTime=")
			require.NotContains(t, req.URL.String(), "sinceTime=1888-06-24")
			require.Contains(t, req.URL.String(), "format=json")
			return mockAPIResponseOK(t), nil
		},
	}
	err = recv.poll(context.Background())
	require.NoError(t, err)

	logs = sink.AllLogs()
	log := logs[0]

	require.Equal(t, 4, log.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().Len())

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestShutdownNoServer(t *testing.T) {
	// test that shutdown without a start does not error or panic
	recv := newReceiver(t, createDefaultConfig().(*Config), consumertest.NewNop())
	require.NoError(t, recv.Shutdown(context.Background()))
}

func newReceiver(t *testing.T, cfg *Config, c consumer.Logs) *proofpointLogsReceiver {
	r, err := newProofpointLogsReceiver(cfg, zap.NewNop(), c)
	require.NoError(t, err)
	return r
}

func mockAPIResponseOK(t *testing.T) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(jsonFileAsString(t, "testdata/proofpointpayload.json"))),
	}
}

func mockAPIResponse429() *http.Response {
	return &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body:       io.NopCloser(strings.NewReader("{}")),
	}
}

func jsonFileAsString(t *testing.T, filePath string) string {
	jsonBytes, err := os.ReadFile(filePath)
	require.NoError(t, err)
	return string(jsonBytes)
}
