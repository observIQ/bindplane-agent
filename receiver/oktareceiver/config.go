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

// Package oktareceiver provides a receiver that receives telemetry from an Okta domain.
package oktareceiver // import "github.com/observiq/bindplane-agent/receiver/oktareceiver"

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	defaultPollInterval = time.Minute

	// OktaTimeFormat ISO 8601 Format
	OktaTimeFormat = "2006-01-02T15:04:05Z"
)

// Config defines the configuration for an Okta receiver
type Config struct {
	// Domain Okta Domain (no https://  -  ex: observiq.okta.com)
	Domain string `mapstructure:"okta_domain"`

	// APIToken Okta API Token
	APIToken string `mapstructure:"api_token"`

	// PollInterval The interval at which the Okta API is scanned for Logs
	// Must be 1s or greater
	PollInterval time.Duration `mapstructure:"poll_interval"`

	// StartTime UTC Timestamp following format specified in OktaTimeFormat
	// Must be within the past 180 days and not in the future
	StartTime string `mapstructure:"start_time"`
}

var (
	errNoDomain            = errors.New("okta_domain must be specified")
	errInvalidDomain       = errors.New("invalid okta_domain, do not include https://")
	errNoAPIToken          = errors.New("api_token must be specified")
	errInvalidPollInterval = errors.New("invalid poll_interval, it must be a duration greater than one second")
)

// Validate ensures an Okta receiver config is correct
func (c *Config) Validate() error {
	if c.Domain == "" {
		return errNoDomain
	}

	if strings.Contains(c.Domain, "https://") || strings.Contains(c.Domain, "http://") {
		return errInvalidDomain
	}

	if c.APIToken == "" {
		return errNoAPIToken
	}

	if c.PollInterval != 0 && c.PollInterval < time.Second {
		return errInvalidPollInterval
	}

	err := validateStartTime(c.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start_time: %w", err)
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
		return errors.New("invalid timestamp: must be in the format YYYY-MM-DDTHH:MM:SS")
	}

	nowUTC := time.Now().UTC()
	time180DaysAgo := nowUTC.AddDate(0, 0, -180)

	// Check if the time is in the valid range the past 180 days and before now
	if parsedTime.Before(time180DaysAgo) || parsedTime.After(nowUTC) {
		return errors.New("invalid timestamp: must be within the past 180 days and not in the future")
	}

	return nil
}
