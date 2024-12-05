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

package googlemanagedprometheusexporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/observiq/bindplane-agent/internal/version"
	gmp "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlemanagedprometheusexporter"
	"go.opentelemetry.io/collector/component"
	"google.golang.org/api/option"
)

const (
	defaultUserAgent = "StanzaLogAgent"
)

// Config is the config the google managed prometheus exporter
type Config struct {
	Credentials     string      `mapstructure:"credentials"`
	CredentialsFile string      `mapstructure:"credentials_file"`
	GMPConfig       *gmp.Config `mapstructure:",squash"`
}

// Validate validates the config
func (c *Config) Validate() error {
	return c.GMPConfig.Validate()
}

// setClientOptions sets the client options used by the GCP config
func (c *Config) setClientOptions() {
	c.GMPConfig.MetricConfig.ClientConfig.GetClientOptions = c.getClientOptions
}

// getClientOptions returns the client options used by the exporter
func (c *Config) getClientOptions() []option.ClientOption {
	opts := []option.ClientOption{}

	switch {
	case c.Credentials != "":
		opts = append(opts, option.WithCredentialsJSON([]byte(c.Credentials)))
	case c.CredentialsFile != "":
		opts = append(opts, option.WithCredentialsFile(c.CredentialsFile))
	}

	return opts
}

// setProject sets the project id from credentials if not already set
func (c *Config) setProject() error {
	if c.GMPConfig.ProjectID != "" {
		return nil
	}

	switch {
	case c.Credentials != "":
		return c.updateProjectFromJSON([]byte(c.Credentials))
	case c.CredentialsFile != "":
		return c.updateProjectFromFile(c.CredentialsFile)
	default:
		return nil
	}
}

func (c *Config) updateProjectFromJSON(jsonBytes []byte) error {
	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		return fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	value, ok := jsonMap["project_id"]
	if !ok {
		return errors.New("project id does not exist")
	}

	strValue, ok := value.(string)
	if !ok {
		return errors.New("project id is not a string")
	}

	c.GMPConfig.ProjectID = strValue
	return nil
}

func (c *Config) updateProjectFromFile(fileName string) error {
	jsonBytes, err := os.ReadFile(filepath.Clean(fileName))
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return c.updateProjectFromJSON(jsonBytes)
}

// createDefaultConfig creates the default config for the exporter
func createDefaultConfig() func() component.Config {
	collectorVersion := version.Version()
	return func() component.Config {
		return &Config{
			GMPConfig: createDefaultGCPConfig(collectorVersion),
		}
	}
}

// createDefaultGCPConfig creates a default GCP config
func createDefaultGCPConfig(collectorVersion string) *gmp.Config {
	factory := gmp.NewFactory()
	config := factory.CreateDefaultConfig().(*gmp.Config)
	config.UserAgent = fmt.Sprintf("%s/%s", defaultUserAgent, collectorVersion)

	return config
}
