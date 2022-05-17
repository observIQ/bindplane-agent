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

package observiq

import (
	"testing"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewAgentConfigManager(t *testing.T) {
	logger := zap.NewNop().Sugar()

	expected := &AgentConfigManager{
		configMap:  make(map[string]string),
		validators: make(map[string]opamp.ValidatorFunc),
		logger:     logger.Named("config manager"),
	}

	actual := NewAgentConfigManager(logger)
	require.Equal(t, expected, actual)
}

func TestAddConfig(t *testing.T) {
	manager := NewAgentConfigManager(zap.NewNop().Sugar())

	configName := "config.json"
	cfgPath := "path/to/config.json"

	manager.AddConfig(configName, cfgPath, opamp.NoopValidator)
	require.Equal(t, cfgPath, manager.configMap[configName])
	// require package cannot Equal on function pointers so we just assert that a validator exists
	require.NotNil(t, manager.validators[configName])
}
