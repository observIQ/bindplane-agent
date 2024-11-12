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

package splunksearchapireceiver

import (
	"errors"
	"strings"
	"time"

	"go.opentelemetry.io/collector/config/confighttp"
)

// Config struct to represent the configuration for the Splunk Search API receiver
type Config struct {
	confighttp.ClientConfig `mapstructure:",squash"`
	Username                string   `mapstructure:"splunk_username"`
	Password                string   `mapstructure:"splunk_password"`
	Searches                []Search `mapstructure:"searches"`
	EventBatchSize          int      `mapstructure:"event_batch_size"`
}

// Search struct to represent a Splunk search
type Search struct {
	Query        string `mapstructure:"query"`
	EarliestTime string `mapstructure:"earliest_time"`
	LatestTime   string `mapstructure:"latest_time"`
	Limit        int    `mapstructure:"limit"`
}

// Validate validates the Splunk Search API receiver configuration
func (cfg *Config) Validate() error {
	if cfg.Endpoint == "" {
		return errors.New("missing Splunk server endpoint")
	}
	if cfg.Username == "" {
		return errors.New("missing Splunk username")
	}
	if cfg.Password == "" {
		return errors.New("missing Splunk password")
	}
	if len(cfg.Searches) == 0 {
		return errors.New("at least one search must be provided")
	}

	for _, search := range cfg.Searches {
		if search.Query == "" {
			return errors.New("missing query in search")
		}

		// query implicitly starts with "search" command
		if !strings.HasPrefix(search.Query, "search ") {
			search.Query = "search " + search.Query
		}

		if strings.Contains(search.Query, "|") {
			return errors.New("command chaining is not supported for queries")
		}

		if search.EarliestTime == "" {
			return errors.New("missing earliest_time in search")
		}
		if search.LatestTime == "" {
			return errors.New("missing latest_time in search")
		}

		// parse time strings to time.Time
		earliestTime, err := time.Parse(time.RFC3339, search.EarliestTime)
		if err != nil {
			return errors.New("earliest_time failed to be parsed as RFC3339")
		}

		latestTime, err := time.Parse(time.RFC3339, search.LatestTime)
		if err != nil {
			return errors.New("latest_time failed to be parsed as RFC3339")
		}

		if earliestTime.UTC().After(latestTime.UTC()) {
			return errors.New("earliest_time must be earlier than latest_time")
		}
		if earliestTime.UTC().After(time.Now().UTC()) {
			return errors.New("earliest_time must be earlier than current time")
		}
		if latestTime.UTC().After(time.Now().UTC()) {
			return errors.New("latest_time must be earlier than current time")
		}
	}
	return nil
}
