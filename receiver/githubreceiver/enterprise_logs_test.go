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

package githubreceiver

import (
	"context"
	"io"
	"net/http"
	"os"
	"strconv"
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
	recv, err := newGitHubLogsReceiver(cfg, zap.NewNop(), consumertest.NewNop())
	require.NoError(t, err)

	err = recv.Start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestShutdownNoServer(t *testing.T) {
	// test that shutdown without a start does not error or panic
	recv := newReceiver(t, createDefaultConfig().(*Config), consumertest.NewNop())
	require.NoError(t, recv.Shutdown(context.Background()))
}

func TestPollBasic(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	mockAccessToken := "AccessToken"
	cfg.AccessToken = mockAccessToken
	cfg.Name = "justin-voss-observiq"
	sink := &consumertest.LogsSink{}
	recv, err := newGitHubLogsReceiver(cfg, zap.NewNop(), sink)
	require.NoError(t, err)
	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			require.Contains(t, req.URL.String(), "per_page="+strconv.Itoa(gitHubMaxLimit))
			return mockAPIResponseOKBasic(t), nil
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

	require.Equal(t, 20, log.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().Len())
	require.Equal(t, mockNextURL, recv.nextURL)
	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestPollEmptyResponse(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	mockAccessToken := "AccessToken"
	cfg.AccessToken = mockAccessToken
	cfg.Name = "justin-voss-observiq"
	sink := &consumertest.LogsSink{}
	recv, err := newGitHubLogsReceiver(cfg, zap.NewNop(), sink)
	require.NoError(t, err)
	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			require.Contains(t, req.URL.String(), "per_page="+strconv.Itoa(gitHubMaxLimit))
			return mockAPIResponseEmpty(t), nil
		},
	}
	err = recv.poll(context.Background())
	require.NoError(t, err)
	logs := sink.AllLogs()
	require.Empty(t, logs)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}
func TestPollTooManyRequests(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	sink := &consumertest.LogsSink{}
	recv, err := newGitHubLogsReceiver(cfg, zap.NewNop(), sink)
	require.NoError(t, err)

	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			recv.client = &mockHTTPClient{
				mockDo: func(req *http.Request) (*http.Response, error) {
					require.Contains(t, req.URL.String(), "per_page="+strconv.Itoa(gitHubMaxLimit))
					return mockAPIResponseTooManyRequests(), nil
				},
			}
			return mockAPIResponseOK100Logs(t), nil
		},
	}

	err = recv.poll(context.Background())
	require.NoError(t, err)

	logs := sink.AllLogs()

	require.Equal(t, 100, logs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().Len())
	require.Equal(t, mockNextURL, recv.nextURL)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}
func TestPollOverflow(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	sink := &consumertest.LogsSink{}
	recv, err := newGitHubLogsReceiver(cfg, zap.NewNop(), sink)
	require.NoError(t, err)

	recv.client = &mockHTTPClient{
		mockDo: func(req *http.Request) (*http.Response, error) {
			recv.client = &mockHTTPClient{
				mockDo: func(req *http.Request) (*http.Response, error) {
					require.Contains(t, req.URL.String(), "per_page="+strconv.Itoa(gitHubMaxLimit))
					return mockAPIResponseOKBasic(t), nil
				},
			}
			return mockAPIResponseOK100Logs(t), nil
		},
	}

	err = recv.poll(context.Background())
	require.NoError(t, err)

	logs := sink.AllLogs()

	require.Equal(t, 120, logs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().Len())
	require.Equal(t, mockNextURL, recv.nextURL)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
}


func newReceiver(t *testing.T, cfg *Config, c consumer.Logs) *gitHubLogsReceiver {
	r, err := newGitHubLogsReceiver(cfg, zap.NewNop(), c)
	require.NoError(t, err)
	return r
}

func jsonFileAsString(t *testing.T, filePath string) string {
	jsonBytes, err := os.ReadFile(filePath)
	require.NoError(t, err)
	return string(jsonBytes)
}

func mockAPIResponseOKBasic(t *testing.T) *http.Response {
	mockRes := &http.Response{}
	mockRes.StatusCode = http.StatusOK
	mockRes.Header = http.Header{}
	mockRes.Header.Add("Link", mockLinkHeaderNext)
	mockRes.Header.Add("Link", mockLinkHeaderLast)
	mockRes.Body = io.NopCloser(strings.NewReader(jsonFileAsString(t, "testdata/gitHubResponseBasicEnterprise.json")))
	return mockRes
}


func mockAPIResponseEmpty(t *testing.T) *http.Response {
	mockRes := &http.Response{}
	mockRes.StatusCode = http.StatusOK
	mockRes.Body = io.NopCloser(strings.NewReader(jsonFileAsString(t, "testdata/gitHubResponseEmpty.json")))
	return mockRes
}

func mockAPIResponseOK100Logs(t *testing.T) *http.Response {
	mockRes := &http.Response{}
	mockRes.StatusCode = http.StatusOK
	mockRes.Header = http.Header{}
	mockRes.Header.Add("Link", mockLinkHeaderNext)
	mockRes.Header.Add("Link", mockLinkHeaderLast)
	mockRes.Body = io.NopCloser(strings.NewReader(jsonFileAsString(t, "testdata/gitHub100ResponseEnterprise.json")))
	return mockRes
}


func mockAPIResponseTooManyRequests() *http.Response {
	return &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body:       io.NopCloser(strings.NewReader("{}")),
	}
}

var (
	mockNextURL        = `https://api.github.com/enterprises/justin-voss-observiq/audit-log?per_page=100&after=MTcyNTQ4MTg2NDI2MXxjZmxuckFiZ1lUT0lUdFdoUi1GUVl3&before=`
	mockLinkHeaderNext = `<https://api.github.com/enterprises/justin-voss-observiq/audit-log?per_page=100&after=MTcyNTQ4MTg2NDI2MXxjZmxuckFiZ1lUT0lUdFdoUi1GUVl3&before=>; rel="next"`
	mockLinkHeaderLast = `<https://api.github.com/enterprises/justin-voss-observiq/audit-log?per_page=100&after=MTcyNTQ4MTg2NDI2MXxjZmxuckFiZ1lUT0lUdFdoUi1GUVl3&before=>; rel="last"`
)
