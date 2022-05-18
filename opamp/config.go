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

// Package opamp contains configurations and protocol implementations to handle OpAmp communication.
package opamp

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	// errPrefixReadFile for error when reading config file
	errPrefixReadFile = "failed to read OpAmp config file"

	// errPrefixParse for error when parsing config
	errPrefixParse = "failed to parse OpAmp config"
)

// Config contains the configuration for the collector to communicate with an OpAmp enabled platform.
type Config struct {
	Endpoint  string  `yaml:"endpoint"`
	SecretKey *string `yaml:"secret_key,omitempty"`
	AgentID   string  `yaml:"agent_id"`

	// Updatable fields
	Labels    *string `yaml:"labels,omitempty"`
	AgentName *string `yaml:"agent_name,omitempty"`
}

// ParseConfig given a configuration file location will parse the config
func ParseConfig(configLocation string) (*Config, error) {
	configPath := filepath.Clean(configLocation)

	// Read in config file contents
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errPrefixReadFile, err)
	}

	// Parse the config
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("%s: %w", errPrefixParse, err)
	}

	return &config, nil
}

// GetSecretKey returns secret key if set else returns empty string
func (c Config) GetSecretKey() string {
	if c.SecretKey == nil {
		return ""
	}

	return *c.SecretKey
}

// CmpUpdatableFields compares updatable fields for equality
func (c Config) CmpUpdatableFields(o Config) (equal bool) {
	if !cmpStringPtr(c.AgentName, o.AgentName) {
		return false
	}

	return cmpStringPtr(c.Labels, o.Labels)
}

func cmpStringPtr(p1, p2 *string) bool {
	switch {
	case p1 == nil && p2 == nil:
		return true
	case p1 == nil && p2 != nil:
		fallthrough
	case p1 != nil && p2 == nil:
		fallthrough
	case *p1 != *p2:
		return false
	default:
		return true
	}
}
