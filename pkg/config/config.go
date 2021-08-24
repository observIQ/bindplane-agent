package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config is the raw config of a collector.
type Config struct {
	Receivers  map[string]map[string]interface{} `yaml:"receivers" mapstructure:"receivers"`
	Processors map[string]map[string]interface{} `yaml:"processors" mapstructure:"processors"`
	Exporters  map[string]map[string]interface{} `yaml:"exporters" mapstructure:"exporters"`
	Extensions map[string]map[string]interface{} `yaml:"extensions" mapstructure:"extensions"`
	Service    Service                           `yaml:"service" mapstructure:"service"`
}

// Service is the raw service config of a collector.
type Service struct {
	Extensions []string            `yaml:"extensions" mapstructure:"extensions"`
	Pipelines  map[string]Pipeline `yaml:"pipelines" mapstructure:"pipelines"`
}

// Pipeline is a raw pipeline config.
type Pipeline struct {
	Receivers  []string `yaml:"receivers" mapstructure:"receivers"`
	Processors []string `yaml:"processors" mapstructure:"processors"`
	Exporters  []string `yaml:"exporters" mapstructure:"exporters"`
}

// Read will read a config from the supplied file path.
func Read(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	config := Config{}
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &config, nil
}
