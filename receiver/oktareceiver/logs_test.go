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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.uber.org/zap"
)

func TestStartShutdown(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	recv, err := newOktaLogsReceiver(cfg, zap.NewNop(), consumertest.NewNop())
	require.NoError(t, err)

	err = recv.Start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err)

	err = recv.Shutdown(context.Background())
	require.NoError(t, err)
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
	recv := newReceiver(t, &Config{
		Domain:       "example.okta.com",
		ApiToken:     "12345",
		PollInterval: time.Minute,
	}, consumertest.NewNop())

	require.NoError(t, recv.Shutdown(context.Background()))
}

func newReceiver(t *testing.T, cfg *Config, c consumer.Logs) *oktaLogsReceiver {
	r, err := newOktaLogsReceiver(cfg, zap.NewNop(), c)
	require.NoError(t, err)
	return r
}
