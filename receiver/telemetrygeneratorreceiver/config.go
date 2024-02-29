// Copyright observIQ, Inc.
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

// Package telemetrygeneratorreceiver generates telemetry for testing purposes
package telemetrygeneratorreceiver //import "github.com/observiq/bindplane-agent/receiver/telemetrygeneratorreceiver"

import (
	"errors"
	"fmt"
)

// Config is the configuration for the telemetry generator receiver
type Config struct {
	PayloadsPerSecond int               `mapstructure:"payloads_per_second"`
	Generators        []GeneratorConfig `mapstructure:"generators"`
}

// GeneratorConfig is the configuration for a single generator
type GeneratorConfig struct {
	// Type of generator to use, either "logs", "host_metrics", or "windows_events"
	Type generatorType `mapstructure:"type"`

	// ResourceAttributes are additional key-value pairs to add to the resource attributes of telemetry.
	ResourceAttributes map[string]string `mapstructure:"resource_attributes"`

	// Attributes are Additional key-value pairs to add to the telemetry attributes
	Attributes map[string]any `mapstructure:"attributes"`

	// AdditionalConfig are any additional config that a generator might need.
	AdditionalConfig map[string]any `mapstructure:"additional_config"`
}

// Validate validates the config
func (c *Config) Validate() error {

	if c.PayloadsPerSecond < 1 {
		return errors.New("payloads_per_second must be at least 1")
	}

	for _, generator := range c.Generators {
		if err := generator.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the generator config
func (g *GeneratorConfig) Validate() error {
	if g.Type == "" {
		return errors.New("type must be set")
	}

	switch g.Type {
	case generatorTypeLogs:
		return validateLogGeneratorConfig(g)
	case generatorTypeHostMetrics:
		return validateHostMetricsGeneratorConfig(g)
	case generatorTypeWindowsEvents:
		return validateWindowsEventsGeneratorConfig(g)

	default:
		return fmt.Errorf("invalid generator type: %s", g.Type)
	}
}

func validateLogGeneratorConfig(g *GeneratorConfig) error {
	// severity and body validation
	if body, ok := g.AdditionalConfig["body"]; ok {
		// check if body is a valid string or map
		// if not, return an error
		_, ok := body.(string)
		if !ok {
			return errors.New("body must be a string")
		}
	}

	// if severity is set, it must be a valid severity
	if severity, ok := g.AdditionalConfig["severity"]; ok {
		if _, ok := severity.(int); !ok {
			return errors.New("severity must be an integer")
		}
	}
	return nil
}

func validateHostMetricsGeneratorConfig(_ *GeneratorConfig) error {
	return nil
}

func validateWindowsEventsGeneratorConfig(_ *GeneratorConfig) error {
	return nil
}
