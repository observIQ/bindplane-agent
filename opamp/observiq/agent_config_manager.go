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
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
)

const (
	// CollectorConfigName is the key of the collector config in OpAmp
	CollectorConfigName = "collector.yaml"
	// ManagerConfigName is the key of the manager config in OpAmp
	ManagerConfigName = "manager.yaml"
	// LoggingConfigName is the key of the logging config in OpAmp
	LoggingConfigName = "logging.yaml"
)

// acceptableConfigs is a lookup of configs that are able to be written/updated
var acceptableConfigs = map[string]struct{}{
	CollectorConfigName: {},
	ManagerConfigName:   {},
	LoggingConfigName:   {},
}

// Enforce interface
var _ opamp.ConfigManager = (*AgentConfigManager)(nil)

// AgentConfigManager keeps track of active configs for the agent
type AgentConfigManager struct {
	configMap map[string]*opamp.ManagedConfig
	logger    *zap.SugaredLogger
}

// NewAgentConfigManager creates a new AgentConfigManager
func NewAgentConfigManager(defaultLogger *zap.SugaredLogger) *AgentConfigManager {
	return &AgentConfigManager{
		configMap: make(map[string]*opamp.ManagedConfig),
		logger:    defaultLogger.Named("config manager"),
	}
}

// AddConfig adds a config to be tracked by the config manager.
// If the config already is tracked it'll be overwritten with the new managed config
func (a *AgentConfigManager) AddConfig(configName string, managedConfig *opamp.ManagedConfig) {
	a.configMap[configName] = managedConfig
}

// ComposeEffectiveConfig reads in all config files and calculates the effective config
func (a *AgentConfigManager) ComposeEffectiveConfig() (*protobufs.EffectiveConfig, error) {
	contentMap := make(map[string]*protobufs.AgentConfigFile, len(a.configMap))

	// Used to track config names to alphabetize later in hash compution
	configNames := make([]string, 0, len(a.configMap))

	for configName, managedConfig := range a.configMap {
		// Read in config file
		cleanPath := filepath.Clean(managedConfig.ConfigPath)
		configContents, err := os.ReadFile(cleanPath)
		if err != nil {
			return nil, fmt.Errorf("error reading config file %s: %w", configName, err)
		}

		// Add to contentMap
		contentMap[configName] = &protobufs.AgentConfigFile{
			Body:        configContents,
			ContentType: opamp.DetermineContentType(managedConfig.ConfigPath),
		}

		configNames = append(configNames, configName)
	}

	return &protobufs.EffectiveConfig{
		ConfigMap: &protobufs.AgentConfigMap{
			ConfigMap: contentMap,
		},
	}, nil
}

// ApplyConfigChanges compares the remoteConfig to the existing and applies changes
func (a *AgentConfigManager) ApplyConfigChanges(remoteConfig *protobufs.AgentRemoteConfig) (effectiveConfig *protobufs.EffectiveConfig, changed bool, err error) {
	effectiveConfig, err = a.ComposeEffectiveConfig()
	if err != nil {
		return nil, false, err
	}

	currentConfigMap := effectiveConfig.GetConfigMap().GetConfigMap()
	remoteConfigMap := remoteConfig.GetConfig().GetConfigMap()

	// No remote config Map
	if remoteConfigMap == nil {
		return
	}

	// loop through all remote configs and compare then with existing configs
	for configName, remoteContents := range remoteConfigMap {
		// For security check the log file we want is acceptable
		if _, ok := acceptableConfigs[configName]; !ok {
			a.logger.Info("Not support config received skipping", zap.String("config", configName))
			continue
		}

		currentContents, ok := currentConfigMap[configName]
		if !ok {
			a.logger.Info("Untracked config found", zap.String("config", configName))
			// We don't current track this config file we should add it
			changed = true

			// Write out the file
			if err := os.WriteFile(configName, remoteContents.GetBody(), 0600); err != nil {
				return nil, false, fmt.Errorf("failed to write new config file %s: %w", configName, err)
			}

			// Track new config
			a.AddConfig(configName, &opamp.ManagedConfig{
				ConfigPath: filepath.Join(".", configName),
				Reload:     opamp.NoopReloadFunc,
			})
			continue
		}

		// Check to see if file contents are the same
		if bytes.Equal(currentContents.GetBody(), remoteContents.GetBody()) {
			continue
		}

		// Update the config file
		managedConfig := a.configMap[configName]
		newContents := remoteContents.GetBody()

		a.logger.Info("Applying changes to config file", zap.String("config", configName))
		configChanged, err := managedConfig.Reload(newContents)
		if err != nil {
			return effectiveConfig, false, fmt.Errorf("failed to reload config: %s: %w", configName, err)
		}

		changed = changed || configChanged
	}

	// If files have changed recompute effective config
	if changed {
		effectiveConfig, err = a.ComposeEffectiveConfig()
		if err != nil {
			return nil, true, err
		}
	}

	return
}
