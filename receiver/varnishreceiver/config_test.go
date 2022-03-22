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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	testDir := t.TempDir()
	testCases := []struct {
		desc                string
		cfg                 Config
		expectedErrContains error
	}{
		{
			desc:                "empty config",
			cfg:                 Config{},
			expectedErrContains: nil,
		},
		{
			desc: "missing exec dir",
			cfg: Config{
				ExecDir: "missing/exec",
			},
			expectedErrContains: errExecDirNotExist,
		},
		{
			desc: "missing working dir",
			cfg: Config{
				WorkingDir: "missing/working",
			},
			expectedErrContains: errWorkingDirNotExist,
		},
		{
			desc: "valid exec and working dir",
			cfg: Config{
				WorkingDir: testDir,
				ExecDir:    testDir,
			},
			expectedErrContains: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualErr := tc.cfg.Validate()
			if tc.expectedErrContains != nil {
				require.Contains(t, actualErr.Error(), tc.expectedErrContains.Error())
			} else {
				require.NoError(t, actualErr)
			}
		})
	}
}
