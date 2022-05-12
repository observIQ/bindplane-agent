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
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/observiq/observiq-otel-collector/opamp"
	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// acceptableConfigs is a lookup of configs that are able to be written/updated
var acceptableConfigs = map[string]struct{}{
	opamp.CollectorConfigName: {},
	opamp.ManagerConfigName:   {},
	opamp.LoggingConfigName:   {},
}

// Enforce interface
var _ opamp.ConfigManager = (*AgentConfigManager)(nil)

// AgentConfigManager keeps track of active configs for the agent
type AgentConfigManager struct {
	configMap  map[string]string
	validators map[string]opamp.ValidatorFunc
	logger     *zap.SugaredLogger
}

// NewAgentConfigManager creates a new AgentConfigManager
func NewAgentConfigManager(defaultLogger *zap.SugaredLogger) *AgentConfigManager {
	return &AgentConfigManager{
		configMap:  make(map[string]string),
		validators: make(map[string]opamp.ValidatorFunc),
		logger:     defaultLogger.Named("config manager"),
	}
}

// AddConfig adds a config to be tracked by the config manager with it's corresponding validator function.
// If the config already is tracked it'll be overwritten with the new configPath
func (a *AgentConfigManager) AddConfig(configName, configPath string, validator opamp.ValidatorFunc) {
	a.configMap[configName] = configPath
	a.validators[configName] = validator
}

// ComposeEffectiveConfig reads in all config files and calculates the effective config
func (a *AgentConfigManager) ComposeEffectiveConfig() (*protobufs.EffectiveConfig, error) {
	contentMap := make(map[string]*protobufs.AgentConfigFile, len(a.configMap))

	// Used to track config names to alphabetize later in hash compution
	configNames := make([]string, 0, len(a.configMap))

	for configName, configPath := range a.configMap {
		// Read in config file
		cleanPath := filepath.Clean(configPath)
		configContents, err := os.ReadFile(cleanPath)
		if err != nil {
			return nil, fmt.Errorf("error reading config file %s: %w", configName, err)
		}

		// Add to contentMap
		contentMap[configName] = &protobufs.AgentConfigFile{
			Body:        configContents,
			ContentType: opamp.DetermineContentType(configPath),
		}

		configNames = append(configNames, configName)
	}

	effectiveConfig := &protobufs.EffectiveConfig{
		ConfigMap: &protobufs.AgentConfigMap{
			ConfigMap: contentMap,
		},
	}

	// Compute Hash
	client.CalcHashEffectiveConfig(effectiveConfig)
	return effectiveConfig, nil
}

// ApplyConfigChanges compares the remoteConfig to the existing and applies changes
func (a *AgentConfigManager) ApplyConfigChanges(remoteConfig *protobufs.AgentRemoteConfig) (effectiveConfig *protobufs.EffectiveConfig, changed bool, err error) {
	effectiveConfig, err = a.ComposeEffectiveConfig()
	if err != nil {
		return nil, false, err
	}

	// No remote config
	if remoteConfig == nil {
		return
	}

	// No config changes return current effective config
	if bytes.Equal(remoteConfig.GetConfigHash(), effectiveConfig.GetHash()) {
		return
	}

	currentConfigMap := effectiveConfig.GetConfigMap().GetConfigMap()
	remoteConfigMap := remoteConfig.GetConfig().GetConfigMap()

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
			a.AddConfig(configName, configName, opamp.NoopValidator)
			continue
		}

		// File has not changed
		if bytes.Equal(currentContents.GetBody(), remoteContents.GetBody()) {
			continue
		}

		a.logger.Info("Changed made to config file, updating", zap.String("config", configName))

		// A single file has changed set overall change to true
		changed = true

		// Update the config file
		validator := a.validators[configName]
		configPath := a.configMap[configName]
		newContents := remoteContents.GetBody()

		// See if we're updating the manager config.
		// Updating the manager config could be a security risk so we want to limit what fields are updatable
		switch configPath {
		case opamp.ManagerConfigName:
			if err := updateManagerConfig(configPath, newContents, validator); err != nil {
				return effectiveConfig, false, fmt.Errorf("failed to update config %s: %w", configName, err)
			}
		default:
			if err := updateConfigFile(configName, configPath, newContents, validator); err != nil {
				return effectiveConfig, false, fmt.Errorf("failed to update config %s: %w", configName, err)
			}
		}
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

func computeHash(configNames []string, contentMap map[string]*protobufs.AgentConfigFile) []byte {
	// Sort config names
	sort.Strings(configNames)

	// Compute hash
	h := sha256.New()

	// Add each config file to the hash in alphabetical order
	for _, configName := range configNames {
		configContents := contentMap[configName].Body
		h.Write(configContents)
	}

	return h.Sum(nil)
}

func updateManagerConfig(configPath string, contents []byte, validator opamp.ValidatorFunc) error {
	// Unmarshal config and only pull fields out that are allowed to be updated.
	var newConfig opamp.Config
	if err := yaml.Unmarshal(contents, &newConfig); err != nil {
		return fmt.Errorf("failed to validate config %s", opamp.ManagerConfigName)
	}

	// Read in existing config file
	currConfig, err := opamp.ParseConfig(configPath)
	if err != nil {
		return err
	}

	// Update allowed fields
	currConfig.AgentName = newConfig.AgentName
	currConfig.Labels = newConfig.Labels

	// Marshal back into bytes
	newContents, err := yaml.Marshal(currConfig)
	if err != nil {
		return fmt.Errorf("failed to reformat manager config: %w", err)
	}

	// Run through update config workflow
	// Note this will rerun the validator which we know works but it keeps file writing in a single place
	return updateConfigFile(opamp.ManagerConfigName, configPath, newContents, validator)
}

func updateConfigFile(configName, configPath string, contents []byte, validator opamp.ValidatorFunc) error {
	// validate file
	if !validator(contents) {
		return fmt.Errorf("failed to validate config %s", configName)
	}

	// Write file
	if err := os.WriteFile(configPath, contents, 0600); err != nil {
		return fmt.Errorf("failed to update config file %s: %w", configName, err)
	}

	return nil
}
