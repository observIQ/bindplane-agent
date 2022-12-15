package alternateprocessor

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/config"
	"go.uber.org/multierr"
)

// Config is the configuration object for the `alternate` processor
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`
	Metrics                  *AlternateRoute `mapstructure:"metrics"`
	Logs                     *AlternateRoute `mapstructure:"logs"`
	Traces                   *AlternateRoute `mapstructure:"traces"`
}

type AlternateRoute struct {
	Rate  *RateTrackerConfig `mapstructure:"rate"`
	Limit float64            `mapstructure:"limit"`
	Route string             `mapstructure:"route"`
}

var (
	errNoRate = errors.New("no rate configuration was specified")
)

// Validate returns whether or not the configuration for the alternate processor is valid
func (c Config) Validate() error {
	var errs error
	if c.Logs != nil {
		errs = multierr.Append(errs, c.Logs.validate("logs"))
	}
	if c.Metrics != nil {
		errs = multierr.Append(errs, c.Metrics.validate("metrics"))
	}
	if c.Traces != nil {
		errs = multierr.Append(errs, c.Traces.validate("traces"))
	}
	return nil
}

func (ar *AlternateRoute) validate(telemetryType string) error {
	if ar.Rate == nil {
		return fmt.Errorf("no rate was specified for telemetry type %s: %w", telemetryType, errNoRate)
	}
	_, err := ar.Rate.Build()
	if err != nil {
		return fmt.Errorf("there was an error parsing the rate for %s: %w", telemetryType, err)
	}

	return nil
}
