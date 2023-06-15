// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package spancountprocessor provides a processor that counts spans and emits the counts as metrics.
package spancountprocessor

import (
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
)

const (
	// defaultMetricName is the default metric name.
	defaultMetricName = "span.count"

	// defaultMetricUnit is the default metric unit.
	defaultMetricUnit = "{spans}"

	// defaultInterval is the default metric interval.
	defaultInterval = time.Minute

	// defaultOTTLMatch is the default match expression.
	defaultOTTLMatch = "true"
)

// Config is the config of the processor.
type Config struct {
	Route          string            `mapstructure:"route"`
	MetricName     string            `mapstructure:"metric_name"`
	MetricUnit     string            `mapstructure:"metric_unit"`
	Interval       time.Duration     `mapstructure:"interval"`
	Match          string            `mapstructure:"match"`
	OTTLMatch      *string           `mapstructure:"ottl_match"`
	Attributes     map[string]string `mapstructure:"attributes"`
	OTTLAttributes map[string]string `mapstructure:"ottl_attributes"`
}

// Validate validates the config, returning an error if the config is invalid
func (c Config) Validate() error {
	if c.Match != "" && c.OTTLMatch != nil {
		return fmt.Errorf("only one of match and ottl_match can be set")
	}

	if c.Attributes != nil && c.OTTLAttributes != nil {
		return fmt.Errorf("only one of attributes and ottl_attributes can be set")
	}

	return nil
}

func (c Config) ottlMatchExpression() string {
	if c.OTTLMatch != nil {
		return *c.OTTLMatch
	}
	return defaultOTTLMatch
}

func (c Config) isOTTL() bool {
	// Use OTTL if the Expr expression is not set
	return c.Match == ""
}

// createDefaultConfig returns the default config for the processor.
func createDefaultConfig() component.Config {
	return &Config{
		MetricName: defaultMetricName,
		MetricUnit: defaultMetricUnit,
		Interval:   defaultInterval,
	}
}
