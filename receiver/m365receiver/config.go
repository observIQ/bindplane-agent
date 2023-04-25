// Copyright  OpenTelemetry Authors
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

package m365receiver // import "github.com/observiq/observiq-otel-collector/receiver/m365receiver"

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/observiq/observiq-otel-collector/receiver/m365receiver/internal/metadata"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	confighttp.HTTPClientSettings           `mapstructure:",squash"`
	MetricsBuilderConfig                    metadata.MetricsBuilderConfig `mapstructure:",squash"`
	Tenant_id                               string                        `mapstructure:"tenant_id"`
	Client_id                               string                        `mapstructure:"client_id"`
	Client_secret                           string                        `mapstructure:"client_secret"`
}

func (c *Config) Validate() error {
	if c.Tenant_id == "" {
		return fmt.Errorf("missing tenant_id; required")
	} else {
		_, err := uuid.Parse(c.Tenant_id)
		if err != nil {
			return fmt.Errorf("tenant_id is invalid; must be a GUID")
		}
	}

	if c.Client_id == "" {
		return fmt.Errorf("missing client_id; required")
	} else {
		_, err := uuid.Parse(c.Client_id)
		if err != nil {
			return fmt.Errorf("client_id is invalid; must be a GUID")
		}
	}

	if c.Client_secret == "" {
		return fmt.Errorf("missing client_secret; required")
	} else {
		re := regexp.MustCompile("^[a-zA-Z0-9-_.~]{1,40}$")
		if !re.MatchString(c.Client_secret) {
			return fmt.Errorf("client_secret is invalid; does not follow correct structure")
		}
	}

	return nil
}
