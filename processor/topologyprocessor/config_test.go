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

package topologyprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigValidate(t *testing.T) {
	t.Run("Default config is valid", func(t *testing.T) {
		err := createDefaultConfig().(*Config).Validate()
		require.NoError(t, err)
	})

	t.Run("Empty configName", func(t *testing.T) {
		cfg := Config{
			Enabled:   true,
			Interval:  defaultInterval,
			AccountID: "myacct",
			OrgID:     "myorg",
		}
		err := cfg.Validate()
		require.Error(t, err)
	})

	t.Run("Empty AccountID", func(t *testing.T) {
		cfg := Config{
			Enabled:    true,
			Interval:   defaultInterval,
			OrgID:      "myorg",
			ConfigName: "myconfig",
		}
		err := cfg.Validate()
		require.Error(t, err)
	})

	t.Run("Empty OrgID", func(t *testing.T) {
		cfg := Config{
			Enabled:    true,
			Interval:   defaultInterval,
			AccountID:  "myacct",
			ConfigName: "myconfig",
		}
		err := cfg.Validate()
		require.Error(t, err)
	})
}
