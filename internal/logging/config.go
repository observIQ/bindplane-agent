package logging

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/observiq/observiq-collector/internal/env"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Collector LoggerConfig `mapstructure:"collector"`
	Manager   LoggerConfig `mapstructure:"manager"`
}

type LoggerConfig struct {
	Level        zapcore.Level `mapstructure:"level"`
	MaxBackups   int           `mapstructure:"max_backups"`
	MaxMegabytes int           `mapstructure:"max_megabytes"`
	MaxDays      int           `mapstructure:"max_days"`
	File         string        `mapstructure:"file"`
}

// DefaultConfig returns the default configuration for logging
func DefaultConfig() (*Config, error) {
	v := getViperInstance()
	config := &Config{}
	err := v.Unmarshal(config)
	return config, err
}

// LoadConfig loads the config from the config path specified through env variables
// 	Fields unspecified will be filled with the values from DefaultConfig
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	v := getViperInstance()
	v.SetConfigFile(path)

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	mapOut := &map[string]interface{}{}
	err = v.Unmarshal(mapOut)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.TextUnmarshallerHookFunc(),
		Result:     config,
	})
	if err != nil {
		return nil, err
	}

	err = dec.Decode(mapOut)
	if err != nil {
		return nil, err
	}

	err = ValidateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func getViperInstance() *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")

	addLoggerConfigDefaults(v, "collector", path.Join(env.HomeDir(), "log", "collector.log"))
	addLoggerConfigDefaults(v, "manager", path.Join(env.HomeDir(), "log", "manager.log"))

	return v
}

func addLoggerConfigDefaults(v *viper.Viper, keyPrefix, filePath string) {
	if !strings.HasSuffix(keyPrefix, ".") {
		keyPrefix += "."
	}

	v.SetDefault(keyPrefix+"level", zap.InfoLevel)
	v.SetDefault(keyPrefix+"max_backups", 3)
	v.SetDefault(keyPrefix+"max_megabytes", 1)
	v.SetDefault(keyPrefix+"max_days", 7)
	v.SetDefault(keyPrefix+"file", filePath)
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
