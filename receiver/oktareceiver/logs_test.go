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

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
)

func TestShutdownNoServer(t *testing.T) {
	// test that shutdown without a start does not error or panic
	recv := newReceiver(t, &Config{
		Domain:   "domain",
		ApiToken: "apitoken",
	}, consumertest.NewNop())

	require.NoError(t, recv.Shutdown(context.Background()))
}

func newReceiver(t *testing.T, cfg *Config, c consumer.Logs) *oktaLogsReceiver {
	r, err := newOktaLogsReceiver(cfg, c)
	require.NoError(t, err)
	return r
}
