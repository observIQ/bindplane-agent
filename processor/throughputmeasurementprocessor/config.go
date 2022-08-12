// Package throughputmeasurementprocessor provides a processor that measure the amount of otlp structures flowing through it
package throughputmeasurementprocessor

import "go.opentelemetry.io/collector/config"

type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`

	Enabled       bool    `mapstructure:"enabled"`
	SamplingRatio float64 `mapstructure:"sampling_ratio"`
}
