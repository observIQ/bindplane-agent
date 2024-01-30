// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snowflakeexporter

import (
	"errors"

	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	defaultLogsSchema    = "logs"
	defaultMetricsSchema = "metrics"
	defaultTracesSchema  = "traces"
	defaultTable         = "data"
)

// Config is the config for the Snowflake exporter
type Config struct {
	exporterhelper.TimeoutSettings `mapstructure:",squash"`
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	configretry.BackOffConfig      `mapstructure:"retry_on_failure"`

	AccountIdentifier string `mapstructure:"account_identifier"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	Database          string `mapstructure:"database"`
	Warehouse         string `mapstructure:"warehouse"`

	Logs    *TelemetryConfig `mapstructure:"logs,omitempty"`
	Metrics *TelemetryConfig `mapstructure:"metrics,omitempty"`
	Traces  *TelemetryConfig `mapstructure:"traces,omitempty"`
}

// TelemetryConfig is a config used by each telemetry type
type TelemetryConfig struct {
	Schema string `mapstructure:"schema,omitempty"`
	Table  string `mapstructure:"table,omitempty"`
}

// Validate ensures the config is correct
func (c *Config) Validate() error {
	if c.AccountIdentifier == "" {
		return errors.New("account_identifier is required")
	}
	if c.Username == "" {
		return errors.New("username is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	if c.Database == "" {
		return errors.New("database is required")
	}
	if c.Warehouse == "" {
		return errors.New("warehouse is required")
	}

	return c.validateTelemetry()
}

// validateTelemetry ensures at least 1 telemetry type is configured and sets default values if needed
func (c *Config) validateTelemetry() error {
	var noTelemetry = true
	if c.Logs != nil {
		noTelemetry = false
		if c.Logs.Schema == "" {
			c.Logs.Schema = defaultLogsSchema
		}
		if c.Logs.Table == "" {
			c.Logs.Table = defaultTable
		}
	}
	if c.Metrics != nil {
		noTelemetry = false
		if c.Metrics.Schema == "" {
			c.Metrics.Schema = defaultMetricsSchema
		}
		if c.Metrics.Table == "" {
			c.Metrics.Table = defaultTable
		}
	}
	if c.Traces != nil {
		noTelemetry = false
		if c.Traces.Schema == "" {
			c.Traces.Schema = defaultTracesSchema
		}
		if c.Traces.Table == "" {
			c.Traces.Table = defaultTable
		}
	}

	if noTelemetry {
		return errors.New("no telemetry type configured for exporter")
	}
	return nil
}
