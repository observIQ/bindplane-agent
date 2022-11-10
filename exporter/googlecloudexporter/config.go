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

package googlecloudexporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector"
	gcp "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.uber.org/multierr"
	"google.golang.org/api/option"
)

const (
	defaultMetricPrefix = "workload.googleapis.com"
	defaultUserAgent    = "observIQ-otel-collector"
)

// Config is the config the google cloud exporter
type Config struct {
	config.ExporterSettings `mapstructure:",squash"`
	Credentials             string                 `mapstructure:"credentials"`
	CredentialsFile         string                 `mapstructure:"credentials_file"`
	AppendHost              bool                   `mapstructure:"append_host"`
	GCPConfig               *gcp.Config            `mapstructure:",squash"`
	BatchConfig             *batchprocessor.Config `mapstructure:"batch"`
}

// Validate validates the config
func (c *Config) Validate() error {
	var err error
	err = multierr.Append(err, c.GCPConfig.Validate())
	err = multierr.Append(err, c.BatchConfig.Validate())
	return err
}

// setClientOptions sets the client options used by the GCP config
func (c *Config) setClientOptions() {
	c.GCPConfig.LogConfig.ClientConfig.GetClientOptions = c.getClientOptions
	c.GCPConfig.MetricConfig.ClientConfig.GetClientOptions = c.getClientOptions
	c.GCPConfig.TraceConfig.ClientConfig.GetClientOptions = c.getClientOptions
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
	if c.GCPConfig.Config.ProjectID != "" {
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

	c.GCPConfig.ProjectID = strValue
	return nil
}

func (c *Config) updateProjectFromFile(fileName string) error {
	jsonBytes, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return c.updateProjectFromJSON(jsonBytes)
}

// createDefaultConfig creates the default config for the exporter
func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		GCPConfig:        createDefaultGCPConfig(),
		BatchConfig:      createDefaultBatchConfig(),
		AppendHost:       true,
	}
}

// createDefaultGCPConfig creates a default GCP config
func createDefaultGCPConfig() *gcp.Config {
	factory := gcp.NewFactory()
	config := factory.CreateDefaultConfig().(*gcp.Config)
	config.RetrySettings.Enabled = false
	config.UserAgent = defaultUserAgent
	config.MetricConfig.Prefix = defaultMetricPrefix
	config.LogConfig.DefaultLogName, _ = os.Hostname()

	// Overwrites the default resource filter to match all resource attributes
	defaultResourceFilter := collector.ResourceFilter{Prefix: ""}
	config.MetricConfig.ResourceFilters = []collector.ResourceFilter{defaultResourceFilter}
	return config
}

// createDefaultBatchConfig creates a default batch config
func createDefaultBatchConfig() *batchprocessor.Config {
	factory := batchprocessor.NewFactory()
	config := factory.CreateDefaultConfig().(*batchprocessor.Config)
	return config
}
