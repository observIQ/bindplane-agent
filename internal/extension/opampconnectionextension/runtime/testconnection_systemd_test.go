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

package runtime

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSystemctlShow(t *testing.T) {
	testCases := []struct {
		name     string
		output   string
		expected systemdUnitInfo
	}{
		{
			name: "loaded with environment and files",
			output: `EnvironmentFiles=/etc/default/observiq-otel-collector (ignore_errors=no)
EnvironmentFiles=/etc/sysconfig/proxy (ignore_errors=yes)
Environment=PATH=/bin:/usr/bin BINDPLANE_COLLECTOR_HOME=/opt/observiq-otel-collector HTTPS_PROXY=http://proxy:3128
LoadState=loaded
`,
			expected: systemdUnitInfo{
				loadState:   "loaded",
				environment: []string{"PATH=/bin:/usr/bin", "BINDPLANE_COLLECTOR_HOME=/opt/observiq-otel-collector", "HTTPS_PROXY=http://proxy:3128"},
				environmentFiles: []systemdEnvironmentFile{
					{path: "/etc/default/observiq-otel-collector"},
					{path: "/etc/sysconfig/proxy", ignoreErrors: true},
				},
			},
		},
		{
			name:   "not found",
			output: "Environment=\nLoadState=not-found\n",
			expected: systemdUnitInfo{
				loadState: "not-found",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := parseSystemctlShow(tc.output)
			require.NoError(t, err)
			require.Equal(t, tc.expected, info)
		})
	}
}

func TestSplitSystemdEnvironment(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    []string
		expectedErr string
	}{
		{
			name:     "empty",
			input:    "",
			expected: nil,
		},
		{
			name:     "plain entries",
			input:    "A=1 B=2",
			expected: []string{"A=1", "B=2"},
		},
		{
			name:     "double quoted entry with escapes",
			input:    `A=1 "NO_PROXY=a b" "C=say \"hi\""`,
			expected: []string{"A=1", "NO_PROXY=a b", `C=say "hi"`},
		},
		{
			name:     "single quoted entry",
			input:    `'A=x y' B=2`,
			expected: []string{"A=x y", "B=2"},
		},
		{
			name:        "unterminated quote",
			input:       `"A=1`,
			expectedErr: "unterminated quote",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entries, err := splitSystemdEnvironment(tc.input)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, entries)
		})
	}
}

func TestParseEnvironmentFile(t *testing.T) {
	data := []byte(`# comment
; another comment

HTTPS_PROXY=http://proxy:3128
  NO_PROXY = "localhost,.internal"
http_proxy='http://proxy:8080'
invalid line
EMPTY=
`)

	require.Equal(t, []string{
		"HTTPS_PROXY=http://proxy:3128",
		"NO_PROXY=localhost,.internal",
		"http_proxy=http://proxy:8080",
		"EMPTY=",
	}, parseEnvironmentFile(data))
}

func TestServiceProxyEnv(t *testing.T) {
	files := map[string]string{
		"/etc/first":  "HTTPS_PROXY=http://first:3128\nOTHER=value\n",
		"/etc/second": "HTTPS_PROXY=http://second:3128\n",
	}
	readFile := func(path string) ([]byte, error) {
		if data, ok := files[path]; ok {
			return []byte(data), nil
		}
		return nil, os.ErrNotExist
	}

	t.Run("files override environment in order", func(t *testing.T) {
		info := systemdUnitInfo{
			environment: []string{"PATH=/bin", "HTTPS_PROXY=http://unit:3128", "NO_PROXY=localhost"},
			environmentFiles: []systemdEnvironmentFile{
				{path: "/etc/first"},
				{path: "/etc/missing", ignoreErrors: true},
				{path: "/etc/second"},
			},
		}

		env, err := serviceProxyEnv(info, readFile)
		require.NoError(t, err)
		require.Equal(t, map[string]string{
			"HTTPS_PROXY": "http://second:3128",
			"NO_PROXY":    "localhost",
		}, env)
	})

	t.Run("missing required file", func(t *testing.T) {
		info := systemdUnitInfo{
			environmentFiles: []systemdEnvironmentFile{{path: "/etc/missing"}},
		}

		_, err := serviceProxyEnv(info, readFile)
		require.ErrorContains(t, err, "failed to read environment file /etc/missing")
		require.True(t, errors.Is(err, os.ErrNotExist))
	})
}

func TestApplyServiceProxyEnv(t *testing.T) {
	clearProxyEnv(t)
	t.Setenv("HTTP_PROXY", "http://shell:3128")
	t.Setenv("NO_PROXY", "shell.internal")

	out := &bytes.Buffer{}
	c := connectionTest{out: out}
	require.NoError(t, c.applyServiceProxyEnv(map[string]string{
		"HTTPS_PROXY": "http://service:3128",
		"NO_PROXY":    "service.internal",
	}))

	require.Equal(t, "http://service:3128", os.Getenv("HTTPS_PROXY"))
	require.Equal(t, "service.internal", os.Getenv("NO_PROXY"))
	_, ok := os.LookupEnv("HTTP_PROXY")
	require.False(t, ok)

	require.Contains(t, out.String(), "Ignoring HTTP_PROXY from the current environment")
	require.NotContains(t, out.String(), "Ignoring NO_PROXY")
}
