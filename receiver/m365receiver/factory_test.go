// Copyright  OpenTelemetry Authors
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

package m365receiver

import (
	"context"
	"testing"
	"time"

	"github.com/observiq/observiq-otel-collector/receiver/m365receiver/internal/metadata"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

func TestType(t *testing.T) {
	factory := NewFactory()
	ft := factory.Type()
	require.EqualValues(t, "m365", ft)
}

func TestCreateMetricsReceiver(t *testing.T) {
	factory := NewFactory()
	test, err := factory.CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		&Config{
			ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
				CollectionInterval: 10 * time.Second,
			},
			HTTPClientSettings: confighttp.HTTPClientSettings{
				Timeout: 10 * time.Second,
			},
			MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
		},
		consumertest.NewNop(),
	)
	require.NoError(t, err)
	require.NotNil(t, test)
}

// func TestValidConfig(t *testing.T) {
// 	factory := NewFactory()
// 	require.NoError(t, component.ValidateConfig(factory.CreateDefaultConfig()))
// }
