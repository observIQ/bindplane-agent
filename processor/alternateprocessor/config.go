package alternateprocessor

import (
	"errors"
	"fmt"
	"time"

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

// AlternateRoute is a name for a config that specifies the next route
type AlternateRoute struct {
	Enabled             bool          `mapstructure:"enabled"`
	Rate                string        `mapstructure:"rate"`
	AggregationInterval time.Duration `mapstructure:"aggregation_interval"`
	Route               string        `mapstructure:"route"`
}

var (
	errNoRate = errors.New("no rate configuration was specified")
)

const (
	defaultAggregationInterval = 10 * time.Second
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
	if ar.Rate == "" {
		return fmt.Errorf("no rate was specified for telemetry type %s: %w", telemetryType, errNoRate)
	}

	r, err := ParseRate(ar.Rate)
	if err != nil {
		return fmt.Errorf("not a valid rate: %w", err)
	}

	if !r.Measure.IsSizeCount() {
		return errors.New("this processor only supports a size throughput rate")
	}

	return nil
}
