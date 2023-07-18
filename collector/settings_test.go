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

package collector

import (
	"context"
	"testing"

	"github.com/observiq/observiq-otel-collector/factories"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

func TestNewSettings(t *testing.T) {
	facts, err := factories.DefaultFactories()
	require.NoError(t, err)

	t.Setenv("FILE", "./test.log")
	configPaths := []string{"./test/valid_with_env_var.yaml"}
	settings, err := NewSettings(configPaths, "0.0.0", nil, facts)
	require.NoError(t, err)

	// Make sure environment variable replacement is working
	provider, err := settings.ConfigProvider.Get(context.Background(), settings.Factories)
	require.NoError(t, err)
	receivcfg := provider.Receivers[component.NewID("filelog")]
	config := receivcfg.(*filelogreceiver.FileLogConfig)
	actualConfVal := config.InputConfig.Include[0]
	require.Equal(t, "./test.log", actualConfVal)

	require.NoError(t, err)
	require.Equal(t, settings.LoggingOptions, []zap.Option(nil))
	require.True(t, settings.DisableGracefulShutdown)
}
