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

	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/utility"
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
	exporterhelper.QueueSettings   `mapstructure:",squash"`
	exporterhelper.RetrySettings   `mapstructure:",squash"`

	AccountIdentifier string `mapstructure:"account_identifier"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	Database          string `mapstructure:"database"`

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

	if err := c.validateTelemetry(); err != nil {
		return err
	}

	if err := c.validateConnection(); err != nil {
		return err
	}
	return nil
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

// validateConnection verifies that the exporter can connect to Snowflake and creates the telemetry schemas if needed
func (c *Config) validateConnection() error {
	// connect to snowflake
	dsn := utility.BuildDSN(
		c.Username,
		c.Password,
		c.AccountIdentifier,
		c.Database,
		"",
	)
	db, err := utility.CreateNewDB(nil, dsn)
	if err != nil {
		return fmt.Errorf("failed to validate connection: %w", err)
	}

	// create schemas if they don't already exist
	if c.Logs != nil {
		if err = utility.CreateSchema(db, c.Logs.Schema); err != nil {
			return err
		}
	}
	if c.Metrics != nil {
		if err = utility.CreateSchema(db, c.Metrics.Schema); err != nil {
			return err
		}
	}
	if c.Traces != nil {
		if err = utility.CreateSchema(db, c.Traces.Schema); err != nil {
			return err
		}
	}

	db.Close()
	return nil
}
