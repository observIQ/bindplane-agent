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

package githubreceiver // import "github.com/observiq/bindplane-agent/receiver/githubreceiver"

import (
	"fmt"
	"time"
)

type WebhookConfig any

// Config defines configuration for GitHub receiver.
type Config struct {
	AccessToken   string        `mapstructure:"access_token"`
	LogType       string        `mapstructure:"log_type"`                // "user", "organization", or "enterprise"
	Name          string        `mapstructure:"name"`                    // The name of the user, organization, or enterprise
	PollInterval  time.Duration `mapstructure:"poll_interval,omitempty"` // Optional
	WebhookConfig WebhookConfig `mapstructure:"webhook,omitempty"`       // Optional
}

// Validate validates the configuration by checking for missing or invalid fields
func (c *Config) Validate() error {
	if c.AccessToken == "" {
		return fmt.Errorf("missing access_token; required")
	}
	if c.LogType == "" {
		return fmt.Errorf("missing log_type; required")
	}
	if c.Name == "" {
		return fmt.Errorf("missing name; required")
	}
	if c.PollInterval == 0 && c.WebhookConfig == nil {
		return fmt.Errorf("must specify either poll_interval or webhook")
	}
	if c.PollInterval < time.Duration(float64(time.Second)*0.72) {
		return fmt.Errorf("invalid poll_interval; must be at least 0.72 seconds")
	}
	return nil
}
