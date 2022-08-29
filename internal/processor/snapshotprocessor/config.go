// Package snapshotprocessor collects metrics, traces, and logs for
package snapshotprocessor

import "go.opentelemetry.io/collector/config"

// Config is the configuration for the processor
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`

	// Enable controls whether snapshots are collected
	Enabled bool `mapstructure:"enabled"`
}

// Validate validates the processor configuration
func (cfg Config) Validate() error {
	return nil
}
