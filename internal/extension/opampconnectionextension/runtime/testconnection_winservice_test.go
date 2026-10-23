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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWindowsServiceProxyEnv(t *testing.T) {
	testCases := []struct {
		name       string
		systemEnv  map[string]string
		serviceEnv []string
		expected   map[string]string
	}{
		{
			name:     "nothing set",
			expected: map[string]string{},
		},
		{
			name:      "system environment only",
			systemEnv: map[string]string{"HTTPS_PROXY": "http://system:3128", "NO_PROXY": "localhost"},
			expected:  map[string]string{"HTTPS_PROXY": "http://system:3128", "NO_PROXY": "localhost"},
		},
		{
			name:       "service environment only",
			serviceEnv: []string{"HTTP_PROXY=http://proxy.corp:8000", "HTTPS_PROXY=http://proxy.corp:8000", "OTHER=value"},
			expected:   map[string]string{"HTTP_PROXY": "http://proxy.corp:8000", "HTTPS_PROXY": "http://proxy.corp:8000"},
		},
		{
			name:       "service environment overrides system environment regardless of case",
			systemEnv:  map[string]string{"HTTPS_PROXY": "http://system:3128", "NO_PROXY": "localhost"},
			serviceEnv: []string{"https_proxy=http://service:3128", "invalid"},
			expected:   map[string]string{"HTTPS_PROXY": "http://service:3128", "NO_PROXY": "localhost"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, windowsServiceProxyEnv(tc.systemEnv, tc.serviceEnv))
		})
	}
}
