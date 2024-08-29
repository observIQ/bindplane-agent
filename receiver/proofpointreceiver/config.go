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

// Package proofpointreceiver provides a receiver that receives telemetry from an Okta domain.
package proofpointreceiver // import "github.com/observiq/bindplane-agent/receiver/proofpointreceiver"

import (
	"errors"
	"time"

	"go.opentelemetry.io/collector/config/configopaque"
)

var (
	defaultPollInterval = 5 * time.Minute
)

// Config defines the configuration for a Proofpoint receiver
type Config struct {
	// Principal
	Principal configopaque.String `mapstructure:"principal"`

	// Secret
	Secret configopaque.String `mapstructure:"secret"`

	// PollInterval The interval at which the Proofpoint API is scanned for Logs
	// Must be at least 1 minute
	PollInterval time.Duration `mapstructure:"poll_interval"`
}

var (
	errNoPrincipal         = errors.New("principal must be specified")
	errNoSecret            = errors.New("secret must be specified")
	errInvalidPollInterval = errors.New("invalid poll_interval, it must be at least 1m")
)

// Validate ensures a Proofpoint receiver config is correct
func (c *Config) Validate() error {
	if c.Principal == "" {
		return errNoPrincipal
	}

	if string(c.Secret) == "" {
		return errNoSecret
	}

	if c.PollInterval < time.Minute {
		return errInvalidPollInterval
	}

	return nil
}
