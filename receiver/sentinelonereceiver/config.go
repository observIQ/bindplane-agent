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

// Package sentinelonereceiver provides a receiver that receives telemetry from SentinelOne.
package sentinelonereceiver // import "github.com/observiq/bindplane-agent/receiver/sentinelonereceiver"

import (
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/config/configopaque"
)

var (
	defaultPollInterval = time.Minute
	defaultAPIs         = []string{
		"activities",
		"agents",
		"threats",
	}
)

// Config defines the configuration for a SentinelOne receiver
type Config struct {
	// BaseURL SentinelOne BaseURL
	BaseURL string `mapstructure:"base_url"`

	// APIToken SentinelOne API Token
	APIToken configopaque.String `mapstructure:"api_token"`

	// PollInterval The interval at which the SentinelOne API is polled for Logs
	// Must be in the range [1 second - 24 hours]
	PollInterval time.Duration `mapstructure:"poll_interval"`

	// APIs Which APIs to poll data from
	// Must be subset of defaultAPIs
	APIs []string `mapstructure:"apis"`
}

var (
	errNoBaseURL           = errors.New("base_url must be specified")
	errNoAPIToken          = errors.New("api_token must be specified")
	errInvalidPollInterval = errors.New("invalid poll_interval, it must be within the range of [1 second - 24 hours]")
	errInvalidAPIs         = errors.New(fmt.Sprintf("invalid apis, must be a subset of %s", defaultAPIs))
)

// Validate ensures a SentinelOne receiver config is correct
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return errNoBaseURL
	}

	if string(c.APIToken) == "" {
		return errNoAPIToken
	}

	if c.PollInterval < time.Second || c.PollInterval > 24*time.Hour {
		return errInvalidPollInterval
	}

	if !sliceIsSubset(c.APIs, defaultAPIs) {
		return errInvalidAPIs
	}

	return nil
}

func sliceIsSubset[T comparable](slice1, slice2 []T) bool {
	// Create a map to track elements in slice2
	elementMap := make(map[T]bool)

	// Populate the map with elements from slice2
	for _, elem := range slice2 {
		elementMap[elem] = true
	}

	// Check if every element in slice1 exists in the map
	for _, elem := range slice1 {
		if !elementMap[elem] {
			return false
		}
	}

	return true
}
