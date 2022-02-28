package pluginreceiver

import (
	"fmt"

	"go.opentelemetry.io/collector/config"
)

// Config is the configuration for a template receiver
type Config struct {
	config.ReceiverSettings
	Template   string                 `mapstructure:"template"`
	Parameters map[string]interface{} `mapstructure:"parameters"`
}

// Validate checks if the Config is valid
func (c *Config) Validate() error {
	return nil
}

// Unmarshal will unmarshal a config.Map into a Config
func (c *Config) Unmarshal(configMap *config.Map) error {
	if configMap == nil || len(configMap.AllKeys()) == 0 {
		return fmt.Errorf("config is empty")
	}

	err := configMap.UnmarshalExact(c)
	if err != nil {
		return err
	}

	return nil
}
