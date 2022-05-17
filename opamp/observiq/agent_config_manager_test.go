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
	"os"
	"path/filepath"
	"testing"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
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

func TestComposeEffectiveConfig(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "File missing from disk error",
			testFunc: func(*testing.T) {
				tmpDir := t.TempDir()
				manager := NewAgentConfigManager(zap.NewNop().Sugar())
				manager.AddConfig("not_real.yaml", filepath.Join(tmpDir, "not_real.yaml"), opamp.NoopValidator)

				effCfg, err := manager.ComposeEffectiveConfig()
				assert.ErrorContains(t, err, "error reading config file")
				assert.Nil(t, effCfg)
			},
		},
		{
			desc: "Multi Config Files",
			testFunc: func(*testing.T) {
				tmpDir := t.TempDir()
				configOne := "one.yaml"
				configOnePath := filepath.Join(tmpDir, configOne)
				configOneContents := []byte(`key: value`)

				configTwo := "two.yaml"
				configTwoPath := filepath.Join(tmpDir, configTwo)
				configTwoContents := []byte(`key2: value2`)

				err := os.WriteFile(configOnePath, configOneContents, 0600)
				assert.NoError(t, err)

				err = os.WriteFile(configTwoPath, configTwoContents, 0600)
				assert.NoError(t, err)

				manager := NewAgentConfigManager(zap.NewNop().Sugar())
				manager.AddConfig(configOne, configOnePath, opamp.NoopValidator)
				manager.AddConfig(configTwo, configTwoPath, opamp.NoopValidator)

				expected := &protobufs.EffectiveConfig{
					ConfigMap: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							configOne: {
								Body:        configOneContents,
								ContentType: opamp.YAMLContentType,
							},
							configTwo: {
								Body:        configTwoContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}

				effCfg, err := manager.ComposeEffectiveConfig()
				assert.NoError(t, err)
				assert.Equal(t, expected, effCfg)

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}

func TestApplyConfigChanges(t *testing.T) {
	testCases := []struct {
		desc     string
		testFunc func(*testing.T)
	}{
		{
			desc: "No remote config To Apply",
			testFunc: func(*testing.T) {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, opamp.ManagerConfigName)
				configContents := []byte(`key: value`)

				err := os.WriteFile(configPath, configContents, 0600)
				assert.NoError(t, err)

				manager := NewAgentConfigManager(zap.NewNop().Sugar())
				manager.AddConfig(opamp.ManagerConfigName, configPath, opamp.NoopValidator)

				remoteConfig := &protobufs.AgentRemoteConfig{
					Config: &protobufs.AgentConfigMap{},
				}

				expectedEffCfg := &protobufs.EffectiveConfig{
					ConfigMap: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        configContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}
				effCfg, changed, err := manager.ApplyConfigChanges(remoteConfig)

				assert.NoError(t, err)
				assert.Equal(t, expectedEffCfg, effCfg)
				assert.False(t, changed)
			},
		},
		{
			desc: "Remote config contains unchanged file",
			testFunc: func(*testing.T) {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, opamp.ManagerConfigName)
				configContents := []byte(`key: value`)

				err := os.WriteFile(configPath, configContents, 0600)
				assert.NoError(t, err)

				manager := NewAgentConfigManager(zap.NewNop().Sugar())
				manager.AddConfig(opamp.ManagerConfigName, configPath, opamp.NoopValidator)

				remoteConfig := &protobufs.AgentRemoteConfig{
					Config: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        configContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}

				expectedEffCfg := &protobufs.EffectiveConfig{
					ConfigMap: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        configContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}
				effCfg, changed, err := manager.ApplyConfigChanges(remoteConfig)

				assert.NoError(t, err)
				assert.Equal(t, expectedEffCfg, effCfg)
				assert.False(t, changed)
			},
		},
		{
			desc: "Remote config contains unknown file",
			testFunc: func(*testing.T) {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, opamp.ManagerConfigName)
				configContents := []byte(`key: value`)

				err := os.WriteFile(configPath, configContents, 0600)
				assert.NoError(t, err)

				manager := NewAgentConfigManager(zap.NewNop().Sugar())
				manager.AddConfig(opamp.ManagerConfigName, configPath, opamp.NoopValidator)

				remoteConfig := &protobufs.AgentRemoteConfig{
					Config: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							"other.yaml": {
								Body:        []byte("other: value"),
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}

				expectedEffCfg := &protobufs.EffectiveConfig{
					ConfigMap: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        configContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}
				effCfg, changed, err := manager.ApplyConfigChanges(remoteConfig)

				assert.NoError(t, err)
				assert.Equal(t, expectedEffCfg, effCfg)
				assert.False(t, changed)
			},
		},
		{
			desc: "Remote config contains untracked known file",
			testFunc: func(*testing.T) {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, opamp.ManagerConfigName)
				configContents := []byte(`key: value`)

				newFileContents := []byte(`logger: value`)

				err := os.WriteFile(configPath, configContents, 0600)
				assert.NoError(t, err)

				manager := NewAgentConfigManager(zap.NewNop().Sugar())
				manager.AddConfig(opamp.ManagerConfigName, configPath, opamp.NoopValidator)

				remoteConfig := &protobufs.AgentRemoteConfig{
					Config: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        configContents,
								ContentType: opamp.YAMLContentType,
							},
							opamp.LoggingConfigName: {
								Body:        newFileContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}

				expectedEffCfg := &protobufs.EffectiveConfig{
					ConfigMap: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        configContents,
								ContentType: opamp.YAMLContentType,
							},
							opamp.LoggingConfigName: {
								Body:        newFileContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}
				effCfg, changed, err := manager.ApplyConfigChanges(remoteConfig)

				assert.NoError(t, err)
				assert.Equal(t, expectedEffCfg, effCfg)
				assert.True(t, changed)
				assert.FileExists(t, filepath.Join(".", opamp.LoggingConfigName))

				err = os.Remove(filepath.Join(".", opamp.LoggingConfigName))
				assert.NoError(t, err)
			},
		},
		{
			desc: "Remote config contains no changes to updatable fields in managerFile",
			testFunc: func(*testing.T) {
				tmpDir := t.TempDir()

				// Setup manager config
				managerCfg := opamp.Config{
					Endpoint: "ws://localhost:1234",
					AgentID:  "6868876d-bf87-466c-bfad-5b6ca2c1f993",
				}
				configContents, err := yaml.Marshal(&managerCfg)
				require.NoError(t, err)

				configPath := filepath.Join(tmpDir, opamp.ManagerConfigName)

				err = os.WriteFile(configPath, configContents, 0600)
				assert.NoError(t, err)

				manager := NewAgentConfigManager(zap.NewNop().Sugar())
				manager.AddConfig(opamp.ManagerConfigName, configPath, opamp.NoopValidator)

				// Write new manager config contents that should not trigger changes
				newManagerCfg := opamp.Config{
					Endpoint: "http://other.host.com",
					AgentID:  "6868876d-bf87-466c-bfad-5b6ca2c1f993",
				}
				newConfigContents, err := yaml.Marshal(&newManagerCfg)
				require.NoError(t, err)

				remoteConfig := &protobufs.AgentRemoteConfig{
					Config: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        newConfigContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}

				expectedEffCfg := &protobufs.EffectiveConfig{
					ConfigMap: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        configContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}
				effCfg, changed, err := manager.ApplyConfigChanges(remoteConfig)

				assert.NoError(t, err)
				assert.Equal(t, expectedEffCfg, effCfg)
				assert.False(t, changed)
			},
		},
		{
			desc: "Remote config contains changes to updatable fields in managerFile",
			testFunc: func(*testing.T) {
				tmpDir := t.TempDir()

				// Setup manager config
				managerCfg := opamp.Config{
					Endpoint: "ws://localhost:1234",
					AgentID:  "6868876d-bf87-466c-bfad-5b6ca2c1f993",
				}
				configContents, err := yaml.Marshal(&managerCfg)
				require.NoError(t, err)

				configPath := filepath.Join(tmpDir, opamp.ManagerConfigName)

				err = os.WriteFile(configPath, configContents, 0600)
				assert.NoError(t, err)

				manager := NewAgentConfigManager(zap.NewNop().Sugar())
				manager.AddConfig(opamp.ManagerConfigName, configPath, opamp.NoopValidator)

				agentName := "name"
				// Write new manager config contents that should trigger changes
				newManagerCfg := opamp.Config{
					Endpoint:  "ws://localhost:1234",
					AgentID:   "6868876d-bf87-466c-bfad-5b6ca2c1f993",
					AgentName: &agentName,
				}
				newConfigContents, err := yaml.Marshal(&newManagerCfg)
				require.NoError(t, err)

				remoteConfig := &protobufs.AgentRemoteConfig{
					Config: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        newConfigContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}

				expectedEffCfg := &protobufs.EffectiveConfig{
					ConfigMap: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.ManagerConfigName: {
								Body:        newConfigContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}
				effCfg, changed, err := manager.ApplyConfigChanges(remoteConfig)

				assert.NoError(t, err)
				assert.Equal(t, expectedEffCfg, effCfg)
				assert.True(t, changed)
			},
		},
		{
			desc: "Remote config contains changed non-manager file",
			testFunc: func(*testing.T) {
				tmpDir := t.TempDir()
				configPath := filepath.Join(tmpDir, opamp.LoggingConfigName)
				configContents := []byte(`key: value`)

				newFileContents := []byte(`logger: value`)

				err := os.WriteFile(configPath, configContents, 0600)
				assert.NoError(t, err)

				manager := NewAgentConfigManager(zap.NewNop().Sugar())
				manager.AddConfig(opamp.LoggingConfigName, configPath, opamp.NoopValidator)

				remoteConfig := &protobufs.AgentRemoteConfig{
					Config: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.LoggingConfigName: {
								Body:        newFileContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}

				expectedEffCfg := &protobufs.EffectiveConfig{
					ConfigMap: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							opamp.LoggingConfigName: {
								Body:        newFileContents,
								ContentType: opamp.YAMLContentType,
							},
						},
					},
				}
				effCfg, changed, err := manager.ApplyConfigChanges(remoteConfig)

				assert.NoError(t, err)
				assert.Equal(t, expectedEffCfg, effCfg)
				assert.True(t, changed)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, tc.testFunc)
	}
}
