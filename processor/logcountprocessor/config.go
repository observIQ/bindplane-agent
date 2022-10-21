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

// Package logcountprocessor provides a processor that counts logs as metrics.
package logcountprocessor

import (
	"time"

	"go.opentelemetry.io/collector/config"
)

const (
	defaultMetricName = "log.count"
	defaultMetricUnit = "{logs}"
	defaultInterval   = time.Minute
	defaultMatch      = "true"
)

// Config is the configuration for the resourceattributetransposer
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`
	Exporter                 config.ComponentID `mapstructure:"exporter"`
	MetricName               string             `mapstructure:"metric_name"`
	MetricUnit               string             `mapstructure:"metric_unit"`
	Interval                 time.Duration      `mapstructure:"interval"`
	Match                    string             `mapstructure:"match"`
	Attributes               map[string]string  `mapstructure:"attributes"`
}

// createDefaultConfig returns the default config for the resourceattributetransposer processor.
func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
		MetricName:        defaultMetricName,
		MetricUnit:        defaultMetricUnit,
		Interval:          defaultInterval,
		Match:             defaultMatch,
	}
}
