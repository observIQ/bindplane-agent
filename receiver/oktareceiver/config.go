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

package oktareceiver // import "github.com/observiq/bindplane-agent/receiver/oktareceiver"

import (
	"errors"
	"fmt"
	"time"
)

// ISO 8601 Format
var OktaTimeFormat = "2006-01-02T15:04:05Z"

// Config defines the configuration for an Okta receiver
type Config struct {
	// Domain Okta Domain (ex: observiq.okta.com)
	Domain string `mapstructure:"okta_domain"`

	// ApiToken Okta Api Token
	ApiToken string `mapstructure:"api_token"`

	// PollInterval The interval at which the Okta API is scanned for Logs
	PollInterval time.Duration `mapstructure:"poll_interval"`

	// StartTime UTC Timestamp following format specified in OktaTimeFormat
	StartTime string `mapstructure:"start_time"`
}

var (
	errNoDomain   = errors.New("a domain must be specified")
	errNoApiToken = errors.New("an api token must be specified")
)

// Validate ensures an Okta receiver config is correct
func (c *Config) Validate() error {
	if c.Domain == "" {
		return errNoDomain
	}

	if c.ApiToken == "" {
		return errNoApiToken
	}

	err := validateStartTime(c.StartTime)
	if err != nil {
		return fmt.Errorf("start_time is invalid: %w", err)
	}

	return nil
}

// validateStartTime validates the passed in timestamp string
// must be within the past 180 days and not in the future
func validateStartTime(startTime string) error {
	if startTime == "" {
		return nil
	}

	parsedTime, err := time.Parse(OktaTimeFormat, startTime)
	if err != nil {
		return errors.New("invalid timestamp, must be in the format YYYY-MM-DDTHH:MM:SS")
	}

	now := time.Now()
	time180DaysAgo := now.AddDate(0, 0, -180)

	// Check if the time is in the valid range the past 180 days and before now
	if parsedTime.Before(time180DaysAgo) || parsedTime.After(now) {
		return errors.New("invalid timestamp, must be within the past 180 days and not in the future")
	}

	return nil
}
