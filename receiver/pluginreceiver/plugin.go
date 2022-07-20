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

package pluginreceiver

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Plugin is a templated pipeline of receivers and processors
type Plugin struct {
	Title       string      `yaml:"title,omitempty"`
	Template    string      `yaml:"template,omitempty"`
	Version     string      `yaml:"version,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Parameters  []Parameter `yaml:"parameters,omitempty"`
}

// LoadPlugin loads a plugin from a file path
func LoadPlugin(path string) (*Plugin, error) {
	cleanPath := filepath.Clean(path)
	bytes, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var plugin Plugin
	if err := yaml.Unmarshal(bytes, &plugin); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin from yaml: %w", err)
	}

	return &plugin, nil
}

// Render renders the plugin's template as a config
func (p *Plugin) Render(values map[string]interface{}) (*RenderedConfig, error) {
	template, err := template.New(p.Title).Parse(p.Template)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin template: %w", err)
	}

	templateValues := p.ApplyDefaults(values)

	var writer bytes.Buffer
	if err := template.Execute(&writer, templateValues); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	renderedCfg, err := NewRenderedConfig(writer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to create rendered config: %w", err)
	}

	return renderedCfg, nil
}

// ApplyDefaults returns a copy of the values map with parameter defaults applied.
// If a value is already present in the map, it supercedes the default.
func (p *Plugin) ApplyDefaults(values map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, parameter := range p.Parameters {
		if parameter.Default == nil {
			continue
		}

		result[parameter.Name] = parameter.Default
	}

	for key, value := range values {
		result[key] = value
	}

	return result
}

// CheckParameters checks the supplied values against the defined parameters of the plugin
func (p *Plugin) CheckParameters(values map[string]interface{}) error {
	if err := p.checkDefined(values); err != nil {
		return fmt.Errorf("definition failure: %w", err)
	}

	if err := p.checkRequired(values); err != nil {
		return fmt.Errorf("required failure: %w", err)
	}

	if err := p.checkType(values); err != nil {
		return fmt.Errorf("type failure: %w", err)
	}

	if err := p.checkSupported(values); err != nil {
		return fmt.Errorf("supported value failure: %w", err)
	}

	return nil
}

// checkDefined checks if any of the supplied values are not defined by the plugin
func (p *Plugin) checkDefined(values map[string]interface{}) error {
	parameterMap := make(map[string]struct{})
	for _, parameter := range p.Parameters {
		parameterMap[parameter.Name] = struct{}{}
	}

	for key := range values {
		if _, ok := parameterMap[key]; !ok {
			return fmt.Errorf("parameter %s is not defined in plugin", key)
		}
	}

	return nil
}

// checkRequired checks if required values are defined
func (p *Plugin) checkRequired(values map[string]interface{}) error {
	for _, parameter := range p.Parameters {
		_, ok := values[parameter.Name]
		if parameter.Required && !ok {
			return fmt.Errorf("parameter %s is missing but required in plugin", parameter.Name)
		}
	}

	return nil
}

// checkType checks if the values match their parameter type
func (p *Plugin) checkType(values map[string]interface{}) error {
	for _, parameter := range p.Parameters {
		value, ok := values[parameter.Name]
		if !ok {
			continue
		}

		switch parameter.Type {
		case stringType:
			if _, ok := value.(string); !ok {
				return fmt.Errorf("parameter %s must be a string", parameter.Name)
			}
		case stringArrayType:
			raw, ok := value.([]interface{})
			if !ok {
				return fmt.Errorf("parameter %s must be a []string", parameter.Name)
			}

			for _, v := range raw {
				if _, ok := v.(string); !ok {
					return fmt.Errorf("parameter %s: expected string, but got %v", parameter.Name, v)
				}
			}
		case intType:
			if _, ok := value.(int); !ok {
				return fmt.Errorf("parameter %s must be an int", parameter.Name)
			}
		case boolType:
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("parameter %s must be a bool", parameter.Name)
			}
		case timezoneType:
			raw, ok := value.(string)
			if !ok {
				return fmt.Errorf("parameter %s must be a string", parameter.Name)
			}
			for tz := range tzlist {
				if raw == tz {
					return nil
				}
			}
			return fmt.Errorf("parameter %s must be a valid timezone", parameter.Name)
		default:
			return fmt.Errorf("unsupported parameter type: %s", parameter.Type)
		}
	}

	return nil
}

// checkSupported checks the values against the plugin's supported values
func (p *Plugin) checkSupported(values map[string]interface{}) error {
OUTER:
	for _, parameter := range p.Parameters {
		if parameter.Supported == nil {
			continue
		}

		value, ok := values[parameter.Name]
		if !ok {
			continue
		}

		for _, v := range parameter.Supported {
			if v == value {
				continue OUTER
			}
		}

		return fmt.Errorf("parameter %s does not match the list of supported values: %v", parameter.Name, parameter.Supported)
	}

	return nil
}

// Parameter is the parameter of plugin
type Parameter struct {
	Name      string        `yaml:"name,omitempty"`
	Type      ParameterType `yaml:"type,omitempty"`
	Default   interface{}   `yaml:"default,omitempty"`
	Supported []interface{} `yaml:"supported,omitempty"`
	Required  bool          `yaml:"required,omitempty"`
}

// ParameterType is the type of a parameter
type ParameterType string

const (
	stringType      ParameterType = "string"
	stringArrayType ParameterType = "[]string"
	boolType        ParameterType = "bool"
	intType         ParameterType = "int"
	timezoneType    ParameterType = "timezone"
)
