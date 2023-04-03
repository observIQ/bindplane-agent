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

package apachedruidreceiver // import "github.com/observiq/observiq-otel-collector/receiver/apachedruidreceiver"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestType(t *testing.T) {
	factory := NewFactory()
	ft := factory.Type()
	require.EqualValues(t, typeStr, ft)
}

func TestCreateMetricsReceiver(t *testing.T) {
	cfg := createDefaultConfig().(*Config)

	_, err := NewFactory().CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		cfg,
		nil,
	)
	require.NoError(t, err)
}

func TestCreateMetricsReceiverBadTLS(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Metrics.TLS = &configtls.TLSServerSetting{
		TLSSetting: configtls.TLSSetting{
			CertFile: "some_cert_file",
			KeyFile:  "some_key_file",
		},
	}

	metrics, err := NewFactory().CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		cfg,
		nil,
	)
	require.Nil(t, metrics)
	require.ErrorContains(t, err, "failed to load TLS config: failed to load TLS cert and key: open some_cert_file: no such file or directory")
}
