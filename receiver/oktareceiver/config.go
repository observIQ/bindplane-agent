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
	"time"
)

// Config defines the configuration for an Okta receiver
type Config struct {
	Domain       string        `mapstructure:"okta_domain"`
	ApiToken     string        `mapstructure:"api_token"`
	PollInterval time.Duration `mapstructure:"poll_interval"`
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

	return nil
}
