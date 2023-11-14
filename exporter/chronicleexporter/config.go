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

package chronicleexporter

import (
	"errors"

	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

var regions = map[string]string{
	"Europe Multi-Region":        "https://europe-backstory.googleapis.com",
	"Frankfurt":                  "https://europe-west3-backstory.googleapis.com",
	"London":                     "http://europe-west2-backstory.googleapis.com",
	"Singapore":                  "https://asia-southeast1-backstory.googleapis.com",
	"Sydney":                     "https://australia-southeast1-backstory.googleapis.com",
	"Tel Aviv":                   "https://me-west1-backstory.googleapis.com",
	"United States Multi-Region": "https://united-states-backstory.googleapis.com",
	"Zurich":                     "https://europe-west6-backstory.googleapis.com",
}

// Config defines configuration for the Chronicle exporter.
type Config struct {
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	// Endpoint is the URL where Chronicle data will be sent.
	Region string `mapstructure:"region"`

	// CredsFilePath is the file path to the Google credentials JSON file.
	CredsFilePath string `mapstructure:"creds_file_path"`

	// Creds are the Google credentials JSON file.
	Creds string `mapstructure:"creds"`

	// LogType is the type of log that will be sent to Chronicle.
	LogType string `mapstructure:"log_type"`

	// RawLogField is the field name that will be used to send raw logs to Chronicle.
	RawLogField string `mapstructure:"raw_log_field"`

	// CustomerID is the customer ID that will be used to send logs to Chronicle.
	CustomerID string `mapstructure:"customer_id"`
}

// Validate checks if the configuration is valid.
func (cfg *Config) Validate() error {
	if cfg.CredsFilePath == "" && cfg.Creds == "" {
		return errors.New("either creds_file_path or creds is required")
	}

	if cfg.LogType == "" {
		return errors.New("log_type is required")
	}

	if cfg.Region == "" {
		return errors.New("region is required")
	}

	if _, ok := regions[cfg.Region]; !ok {
		return errors.New("region is invalid")
	}

	return nil
}
