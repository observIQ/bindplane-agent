package resourceattributetransposerprocessor

import "go.opentelemetry.io/collector/config"

// CopyResourceConfig is a config struct specifying a mapping of a resource attribute to a datapoint attribute
type CopyResourceConfig struct {
	// From is the attribute on the resource to copy from
	From string `mapstructure:"from"`
	// To is the attribute to copy to on the individual data point
	To string `mapstructure:"to"`
}

// Config is the configuration for the resourceattributetransposer
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`
	// Operations is a list of copy operations to perform on each ResourceMetric.
	Operations []CopyResourceConfig `mapstructure:"operations"`
}
