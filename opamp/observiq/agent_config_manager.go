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
	"github.com/open-telemetry/opamp-go/protobufs"
)

// Enforce interface
var _ opamp.ConfigManager = (*AgentConfigManager)(nil)

// AgentConfigManager keeps track of active configs for the agent
type AgentConfigManager struct {
	configMap  map[string]string
	validators map[string]opamp.ValidatorFunc
}

// NewAgentConfigManager creates a new AgentConfigManager
func NewAgentConfigManager() *AgentConfigManager {
	return &AgentConfigManager{
		configMap: make(map[string]string),
	}
}

// AddConfig adds a config to be tracked by the config manager with it's corresponding validator function.
// If the config already is tracked it'll be overwritten with the new configPath
func (a *AgentConfigManager) AddConfig(configName, configPath string, validator opamp.ValidatorFunc) {
	cleanPath := filepath.Clean(configPath)

	a.configMap[configName] = cleanPath
	a.validators[configName] = validator
}

// ComposeEffectiveConfig reads in all config files and calculates the effective config
func (a *AgentConfigManager) ComposeEffectiveConfig() (*protobufs.EffectiveConfig, error) {
	contentMap := make(map[string]*protobufs.AgentConfigFile, len(a.configMap))

	// Used to track config names to alphabetize later in hash compution
	configNames := make([]string, 0, len(a.configMap))

	for configName, configPath := range a.configMap {
		// Read in config file
		configContents, err := os.ReadFile(configPath)
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

	// Compute Hash
	configHash := computeHash(configNames, contentMap)

	return &protobufs.EffectiveConfig{
		Hash: configHash,
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
		currentContents, ok := currentConfigMap[configName]
		if !ok {
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

		// A single file has changed set overall change to true
		changed = true

		validator := a.validators[configName]
		if !validator(remoteContents.GetBody()) {
			return effectiveConfig, false, fmt.Errorf("Failed to validate config %s: %w", configName, err)
		}

		// Write file
		filePath := a.configMap[configName]
		if err := os.WriteFile(filePath, remoteContents.GetBody(), 0600); err != nil {
			return nil, false, fmt.Errorf("failed to update config file %s: %w", configName, err)
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

func writeConfigFile(filePath string, contents []byte) error {
	return os.WriteFile(filePath, contents, 0600)
}
