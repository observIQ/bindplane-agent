// Copyright  observIQ, Inc.
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

// Package topologyprocessor collects metrics, traces, and logs for
package topologyprocessor

import (
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"
)

var defaultOpAMPExtensionID = component.MustNewID("opamp")

const defaultInterval = time.Minute

// Config is the configuration for the processor
type Config struct {
	// Interval is the interval at which this processor sends topology messages to BindPlane
	Interval   time.Duration `mapstructure:"interval"`
	OpAMP      component.ID  `mapstructure:"opamp"`
	ConfigName string        `mapstructure:"configName"`
	OrgID      string        `mapstructure:"orgID"`
	AccountID  string        `mapstructure:"accountID"`
}

// Validate validates the processor configuration
func (cfg Config) Validate() error {
	var emptyID component.ID
	if cfg.OpAMP == emptyID {
		return errors.New("`opamp` must be specified")
	}

	if cfg.ConfigName == "" {
		return errors.New("`configName` must be specified")
	}

	if cfg.OrgID == "" {
		return errors.New("`orgID` must be specified")
	}

	if cfg.AccountID == "" {
		return errors.New("`accountID` must be specified")
	}

	return nil
}
