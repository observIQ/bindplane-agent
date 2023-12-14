// Copyright observIQ, Inc.
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

package chronicleforwarderexporter

import (
	"fmt"

	"github.com/observiq/bindplane-agent/expr"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.uber.org/zap"
)

const (
	// ExportTypeSyslog is the syslog export type.
	ExportTypeSyslog = "syslog"

	// ExportTypeFile is the file export type.
	ExportTypeFile = "file"
)

// Config defines configuration for the Chronicle exporter.
type Config struct {
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	// ExportType is the type of export to use.
	ExportType string `mapstructure:"export_type"`

	// Syslog is the configuration for the connection to the Chronicle forwarder.
	Syslog SyslogConfig `mapstructure:"syslog"`

	// File is the configuration for the connection to the Chronicle forwarder.
	File File `mapstructure:"file"`

	// RawLogField is the field name that will be used to send raw logs to Chronicle.
	RawLogField string `mapstructure:"raw_log_field"`
}

// SyslogConfig defines configuration for the Chronicle forwarder connection.
type SyslogConfig struct {
	// Host is the Chronicle forwarder endpoint to send logs to.
	Host string `mapstructure:"host"`

	// port is the port to send logs to.
	Port int `mapstructure:"port"`

	// Network is the network protocol to use.
	Network string `mapstructure:"network"`

	confighttp.HTTPServerSettings `mapstructure:",squash"`
}

// File defines configuration for sending to.
type File struct {
	// Path is the path to the file to send to Chronicle.
	Path string `mapstructure:"path"`
}

func (cfg *Config) Validate() error {
	if cfg.ExportType != ExportTypeSyslog && cfg.ExportType != ExportTypeFile {
		return fmt.Errorf("export_type must be either 'syslog' or 'file'")
	}

	if cfg.ExportType == ExportTypeSyslog {
		if cfg.Syslog.Host == "" || cfg.Syslog.Port <= 0 || cfg.Syslog.Network == "" {
			return fmt.Errorf("incomplete syslog configuration: host, port, and network are required")
		}
	}

	if cfg.ExportType == ExportTypeFile {
		if cfg.File.Path == "" {
			return fmt.Errorf("file path is required for file export type")
		}
	}

	if cfg.RawLogField != "" {
		_, err := expr.NewOTTLLogRecordExpression(cfg.RawLogField, component.TelemetrySettings{
			Logger: zap.NewNop(),
		})
		if err != nil {
			return fmt.Errorf("raw_log_field is invalid: %s", err)
		}
	}
	return nil
}
