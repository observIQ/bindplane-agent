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

	"go.opentelemetry.io/collector/component"
)

// Config is the configuration for the telemetry generator receiver
type Config struct {
	PayloadsPerSecond int               `mapstructure:"payloads_per_second"`
	Generators        []GeneratorConfig `mapstructure:"generators"`
}

// GeneratorConfig is the configuration for a single generator
type GeneratorConfig struct {
	// Type of generator to use, either "logs", "metrics", or "traces"
	Type component.DataType `mapstructure:"type"`

	// ResourceAttributes are additional key-value pairs to add to the resource attributes of telemetry.
	ResourceAttributes map[string]string `mapstructure:"resource_attributes"`

	// Attributes are Additional key-value pairs to add to the telemetry attributes
	Attributes map[string]string `mapstructure:"attributes"`

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

	if g.Type != component.DataTypeLogs && g.Type != component.DataTypeMetrics && g.Type != component.DataTypeTraces {
		return errors.New("type must be one of logs, metrics, or traces")
	}

	// severity and body validation
	if g.Type == component.DataTypeLogs {
		if body, ok := g.AdditionalConfig["body"]; ok {

			// check if body is a valid string or map
			// if not, return an error
			_, ok := body.(string)
			if !ok {
				_, ok := body.(map[string]any)
				if !ok {
					return errors.New("body must be a string or a map")
				}
			}
		}

		// if severity is set, it must be a valid severity
		if severity, ok := g.AdditionalConfig["severity"]; ok {
			if _, ok := severity.(int); !ok {
				return errors.New("severity must be an integer")
			}
		}
	}

	return nil
}
