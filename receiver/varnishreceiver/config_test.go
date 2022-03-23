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

package varnishreceiver // import "github.com/observiq/observiq-otel-collector/receiver/varnishreceiver"

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	testDir := t.TempDir()
	testCases := []struct {
		desc                string
		cfg                 Config
		expectedErrContains string
	}{
		{
			desc:                "empty config",
			cfg:                 Config{},
			expectedErrContains: "",
		},
		{
			desc: "missing exec dir",
			cfg: Config{
				ExecDir: "missing/exec",
			},
			expectedErrContains: `"exec_dir" does not exists`,
		},
		{
			desc: "missing instance name",
			cfg: Config{
				InstanceName: "missing/instance_name",
			},
			expectedErrContains: `"instance_name" does not exists`,
		},
		{
			desc: "valid exec and instance name",
			cfg: Config{
				InstanceName: testDir,
				ExecDir:      testDir,
			},
			expectedErrContains: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualErr := tc.cfg.Validate()
			if tc.expectedErrContains != "" {
				require.Contains(t, actualErr.Error(), tc.expectedErrContains)
			} else {
				require.NoError(t, actualErr)
			}
		})
	}
}

func TestSetDefaultHostname(t *testing.T) {
	t.Run("set default hostname to empty config instance name", func(t *testing.T) {
		cfg := Config{}
		hostname, err := os.Hostname()
		require.NoError(t, err)
		err = cfg.SetDefaultHostname()
		require.NoError(t, err)
		require.EqualValues(t, cfg.InstanceName, hostname)
	})

	t.Run("reuse existing config instance name", func(t *testing.T) {
		cfg := Config{
			InstanceName: "varnishcache",
		}
		err := cfg.SetDefaultHostname()
		require.NoError(t, err)
		require.EqualValues(t, cfg.InstanceName, "varnishcache")
	})
}
