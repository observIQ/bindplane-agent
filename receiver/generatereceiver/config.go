// Package generatereceiver is a package that generates telemetry.
package generatereceiver

import (
	"time"

	"go.opentelemetry.io/collector/component"
)

// Config is the configuration for the generate receiver.
type Config struct {
	Logs *LogConfig `mapstructure:"logs"`
}

// LogConfig is a struct that contains configuration for generating logs.
type LogConfig struct {
	Interval   time.Duration  `mapstructure:"interval"`
	Resource   map[string]any `mapstructure:"resource"`
	Attributes map[string]any `mapstructure:"attributes"`
	Body       any            `mapstructure:"body"`
}

// createDefaultConfig creates the default configuration for the receiver.
func createDefaultConfig() component.Config {
	return &Config{
		Logs: createDefaultLogConfig(),
	}
}

// createDefaultLogConfig creates the default configuration for the log config.
func createDefaultLogConfig() *LogConfig {
	return &LogConfig{
		Interval:   1 * time.Second,
		Resource:   make(map[string]any),
		Attributes: make(map[string]any),
		Body:       "Hello, world!",
	}
}
