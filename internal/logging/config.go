package logging

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/observiq/observiq-collector/internal/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Collector LoggerConfig `yaml:"collector" mapstructure:"collector"`
	Manager   LoggerConfig `yaml:"manager" mapstructure:"manager"`
}

type LoggerConfig struct {
	Level        zapcore.Level `yaml:"level" mapstructure:"level"`
	MaxBackups   int           `yaml:"max_backups" mapstructure:"max_backups"`
	MaxMegabytes int           `yaml:"max_megabytes" mapstructure:"max_megabytes"`
	MaxDays      int           `yaml:"max_days" mapstructure:"max_days"`
	File         string        `yaml:"file" mapstructure:"file"`
}

// DefaultConfig returns the default configuration for logging
func DefaultConfig() *Config {
	return &Config{
		Collector: LoggerConfig{
			Level:        zap.InfoLevel,
			MaxBackups:   3,
			MaxMegabytes: 1,
			MaxDays:      7,
			File:         path.Join(env.DefaultEnvProvider.LogDir(), "collector.log"),
		},
		Manager: LoggerConfig{
			Level:        zap.InfoLevel,
			MaxBackups:   3,
			MaxMegabytes: 1,
			MaxDays:      7,
			File:         path.Join(env.DefaultEnvProvider.LogDir(), "manager.log"),
		},
	}
}

// LoadConfig loads the config from the config path specified through env variables
// 	Fields unspecified will be filled with the values from DefaultConfig
func LoadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}

	err = ValidateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfig checks that the passed config is a valid configuration
//	If something is incorrect with the config, an error will be returned indicating the issue.
func ValidateConfig(c *Config) error {
	err := validateLoggerConfig(&c.Manager)
	if err != nil {
		return fmt.Errorf("error validating Manager logger config: %w", err)
	}

	err = validateLoggerConfig(&c.Collector)
	if err != nil {
		return fmt.Errorf("error validating Collector logger config: %w", err)
	}

	return nil
}

// validateLoggerConfig validates a LoggerConfig.
//  If something is wrong with the LoggerConfig, an error will be returned indicating the issue.
func validateLoggerConfig(c *LoggerConfig) error {
	if c.File == "" {
		return errors.New("LoggerConfig.File must not be empty")
	}
	if c.MaxBackups < 0 {
		return errors.New("LoggerConfig.MaxBackups must be >= 0")
	}

	if c.MaxDays < 0 {
		return errors.New("LoggerConfig.MaxDays must be >= 0")
	}

	if c.MaxMegabytes < 0 {
		return errors.New("LoggerConfig.MaxMegabytes must be >= 0")
	}

	return nil
}
