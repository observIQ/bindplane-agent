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

package oktareceiver

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/plogtest"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type mockHTTPClient struct {
	mockDo func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.mockDo(req)
}

func TestStartShutdown(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	recv, err := newOktaLogsReceiver(cfg, zap.NewNop(), consumertest.NewNop())
	require.NoError(t, err)

	err = recv.Start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestStartContextDone(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	recv, err := newOktaLogsReceiver(cfg, zap.NewNop(), consumertest.NewNop())
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	err = recv.Start(ctx, componenttest.NewNopHost())
	require.NoError(t, err)

	cancel()
}

func TestStartTimeParse(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.StartTime = "2024-08-12T00:00:00Z"

	recv, err := newOktaLogsReceiver(cfg, zap.NewNop(), consumertest.NewNop())
	require.NoError(t, err)
	require.Equal(t, time.Date(2024, 8, 12, 0, 0, 0, 0, time.UTC), recv.startTime)
}

func TestStartTimeParseError(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.StartTime = "9999999-08-12T00:00:00Z"

	_, err := newOktaLogsReceiver(cfg, zap.NewNop(), consumertest.NewNop())
	require.Error(t, err)
}

func TestShutdownNoServer(t *testing.T) {
	// test that shutdown without a start does not error or panic
	recv := newReceiver(t, createDefaultConfig().(*Config), consumertest.NewNop())
	require.NoError(t, recv.Shutdown(context.Background()))
}

func TestPoll(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	mockDomain := "observiq.okta.com"
	cfg.Domain = mockDomain

	sink := &consumertest.LogsSink{}
	recv, err := newOktaLogsReceiver(cfg, zap.NewNop(), sink)
	require.NoError(t, err)

	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			require.Contains(t, req.URL.String(), "https://observiq.okta.com/api/v1/logs?since")
			return mockAPIResponse200(), nil
		},
	}

	err = recv.poll(context.Background())
	require.NoError(t, err)

	logs := sink.AllLogs()

	log := logs[0]
	oktaDomain, exist := log.ResourceLogs().At(0).Resource().Attributes().Get("okta.domain")
	require.True(t, exist)
	require.Equal(t, mockDomain, oktaDomain.Str())
	expected, err := jsonFileAsPlogs("testdata/plog.json")
	require.NoError(t, err)
	require.NoError(t, plogtest.CompareLogs(expected, log, plogtest.IgnoreObservedTimestamp()))

	require.Equal(t, "https://observiq.okta.com/api/v1/logs?limit=20&after=1627500044869_1", recv.nextURL)

	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			require.Equal(t, "https://observiq.okta.com/api/v1/logs?limit=20&after=1627500044869_1", req.URL.String())
			return mockAPIResponse200(), nil
		},
	}
	err = recv.poll(context.Background())
	require.NoError(t, err)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}

func newReceiver(t *testing.T, cfg *Config, c consumer.Logs) *oktaLogsReceiver {
	r, err := newOktaLogsReceiver(cfg, zap.NewNop(), c)
	require.NoError(t, err)
	return r
}

func jsonFileAsPlogs(filepath string) (plog.Logs, error) {
	expectedFileBytes, err := os.ReadFile(filepath)
	if err != nil {
		return plog.Logs{}, nil
	}
	unmarshaler := &plog.JSONUnmarshaler{}
	return unmarshaler.UnmarshalLogs(expectedFileBytes)
}

func mockAPIResponse200() *http.Response {
	mockRes := &http.Response{}
	mockRes.StatusCode = http.StatusOK
	mockRes.Header = http.Header{}
	mockRes.Header.Add("Link", mockLinkHeaderNext)
	mockRes.Header.Add("Link", mockLinkHeaderSelf)
	mockRes.Body = io.NopCloser(strings.NewReader(jsonFileAsString("testdata/oktaResponse.json")))
	return mockRes
}

func jsonFileAsString(filePath string) string {
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %s", err)
		return ""
	}
	return string(jsonBytes)
}

var (
	mockLinkHeaderNext = `<https://observiq.okta.com/api/v1/logs?limit=20&after=1627500044869_1>; rel="next"`
	mockLinkHeaderSelf = `<https://observiq.okta.com/api/v1/logs?limit=20>; rel="self"`
)
