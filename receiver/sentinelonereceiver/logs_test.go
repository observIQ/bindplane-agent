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

package sentinelonereceiver

import (
	"context"
	"net/http"
	"testing"

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

	recv, err := newSentinelOneLogsReceiver(cfg, zap.NewNop(), consumertest.NewNop())
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

func newReceiver(t *testing.T, cfg *Config, c consumer.Logs) *sentinelOneLogsReceiver {
	r, err := newSentinelOneLogsReceiver(cfg, zap.NewNop(), c)
	require.NoError(t, err)
	return r
}
