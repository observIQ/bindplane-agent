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
	"fmt"

	"github.com/snowflakedb/gosnowflake"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	defaultDatabase      = "bpop"
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

	AccountIdentifier string             `mapstructure:"account_identifier"`
	Username          string             `mapstructure:"username"`
	Password          string             `mapstructure:"password"`
	Warehouse         string             `mapstructure:"warehouse"`
	Role              string             `mapstructure:"role,omitempty"`
	Database          string             `mapstructure:"database,omitempty"`
	Parameters        map[string]*string `mapstructure:"parameters,omitempty"`

	Logs    TelemetryConfig `mapstructure:"logs"`
	Metrics TelemetryConfig `mapstructure:"metrics"`
	Traces  TelemetryConfig `mapstructure:"traces"`

	DSN string
}

// TelemetryConfig is a config used by each telemetry type
type TelemetryConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Schema  string `mapstructure:"schema,omitempty"`
	Table   string `mapstructure:"table,omitempty"`
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
	if c.Warehouse == "" {
		return errors.New("warehouse is required")
	}
	if c.Database == "" {
		c.Database = defaultDatabase
	}

	if err := c.validateTelemetry(); err != nil {
		return err
	}

	sf := gosnowflake.Config{}

	sf.User = c.Username
	sf.Password = c.Password
	sf.Account = c.AccountIdentifier
	if c.Parameters != nil {
		sf.Params = c.Parameters
	}

	dsn, err := gosnowflake.DSN(&sf)
	if err != nil {
		return fmt.Errorf("failed to build DSN: %w", err)
	}

	c.DSN = dsn

	return nil
}

// validateTelemetry ensures at least 1 telemetry type is configured and sets default values if needed
func (c *Config) validateTelemetry() error {
	if !c.Logs.Enabled && !c.Metrics.Enabled && !c.Traces.Enabled {
		return errors.New("no telemetry type configured for exporter")
	}

	if c.Logs.Enabled {
		if c.Logs.Schema == "" {
			c.Logs.Schema = defaultLogsSchema
		}
		if c.Logs.Table == "" {
			c.Logs.Table = defaultTable
		}
	}
	if c.Metrics.Enabled {
		if c.Metrics.Schema == "" {
			c.Metrics.Schema = defaultMetricsSchema
		}
		if c.Metrics.Table == "" {
			c.Metrics.Table = defaultTable
		}
	}
	if c.Traces.Enabled {
		if c.Traces.Schema == "" {
			c.Traces.Schema = defaultTracesSchema
		}
		if c.Traces.Table == "" {
			c.Traces.Table = defaultTable
		}
	}

	return nil
}
