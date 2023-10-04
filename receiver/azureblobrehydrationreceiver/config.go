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

package azureblobrehydrationreceiver //import "github.com/observiq/bindplane-agent/receiver/azureblobrehydrationreceiver"

import (
	"errors" // timeFormat is the format for the starting and end time
	"fmt"
	"time"
)

const timeFormat = "2006-01-02T15:04"

// Config is the configuration for the azure blob rehydration receiver
type Config struct {
	// ConnectionString is the Azure Blob Storage connection key,
	// which can be found in the Azure Blob Storage resource on the Azure Portal. (no default)
	ConnectionString string `mapstructure:"connection_string"`

	// Container is the name of the user created storage container. (no default)
	Container string `mapstructure:"container"`

	// RootFolder is the name of the root folder in path.
	RootFolder string `mapstructure:"root_folder"`

	// StartingTime the UTC timestamp to start rehydration from.
	StartingTime string `mapstructure:"starting_time"`

	// EndingTime the UTC timestamp to rehydrate up until.
	EndingTime string `mapstructure:"ending_time"`

	// DeleteOnRead indicates if a file should be deleted once it has been processed
	DeleteOnRead bool `mapstructure:"delete_on_read"`
}

// Validate validates the config
func (c *Config) Validate() error {
	if c.ConnectionString == "" {
		return errors.New("connection_string is required")
	}

	if c.Container == "" {
		return errors.New("container is required")
	}

	if err := validateTimestamp(c.StartingTime); err != nil {
		return fmt.Errorf("starting_time is invalid: %w", err)
	}

	if err := validateTimestamp(c.EndingTime); err != nil {
		return fmt.Errorf("ending_time is invalid: %w", err)
	}

	return nil
}

// validateTimestamp validates the passed in timestamp string
func validateTimestamp(timestamp string) error {
	if timestamp == "" {
		return errors.New("missing value")
	}

	if _, err := time.Parse(timeFormat, timestamp); err != nil {
		return errors.New("invalid timestamp format must be in the form YYYY-MM-DDTHH:MM")
	}

	return nil
}
