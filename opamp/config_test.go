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

package opamp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	// Keep this outside so it can be referenced as pointer
	secretKeyContents := "b92222ee-a1fc-4bb1-98db-26de3448541b"

	testCases := []struct {
		desc                string
		createFile          bool // Used to call file read failure
		configContents      string
		expectedConfig      *Config
		expectedErrContents *string
	}{
		{
			desc:                "Failed File Read",
			createFile:          false,
			configContents:      "",
			expectedConfig:      nil,
			expectedErrContents: &errPrefixReadFile,
		},
		{
			desc:       "Failed Marshal",
			createFile: true,
			configContents: `
			{
				"endpoint": "localhost:1234"
			}`,
			expectedConfig:      nil,
			expectedErrContents: &errPrefixParse,
		},
		{
			desc:       "Successful Full Parse",
			createFile: true,
			configContents: `
endpoint: localhost:1234
secret_key: b92222ee-a1fc-4bb1-98db-26de3448541b
agent_id: 8321f735-a52c-4f49-aca9-66f9266c5fe5
labels:
  - one
  - two
`,
			expectedConfig: &Config{
				Endpoint:  "localhost:1234",
				SecretKey: &secretKeyContents,
				AgentID:   "8321f735-a52c-4f49-aca9-66f9266c5fe5",
				Labels: []string{
					"one",
					"two",
				},
			},
			expectedErrContents: nil,
		},
		{
			desc:       "Successful Partial Parse",
			createFile: true,
			configContents: `
endpoint: localhost:1234
agent_id: 8321f735-a52c-4f49-aca9-66f9266c5fe5
`,
			expectedConfig: &Config{
				Endpoint:  "localhost:1234",
				SecretKey: nil,
				AgentID:   "8321f735-a52c-4f49-aca9-66f9266c5fe5",
				Labels:    nil,
			},
			expectedErrContents: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "manager.yml")

			// Only create a config file if we're configured too.
			// This exists to trigger a file read failure
			if tc.createFile {
				err := os.WriteFile(configPath, []byte(tc.configContents), os.ModePerm)
				require.NoError(t, err)
			}

			cfg, err := ParseConfig(configPath)
			if tc.expectedErrContents == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, *tc.expectedErrContents)
			}

			require.Equal(t, tc.expectedConfig, cfg)
		})
	}
}
