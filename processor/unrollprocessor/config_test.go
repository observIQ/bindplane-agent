// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unrollprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {

	testCases := []struct {
		desc        string
		cfg         *Config
		expectedErr string
	}{
		{
			desc: "valid config",
			cfg:  createDefaultConfig().(*Config),
		},
		{
			desc: "config without body field",
			cfg: &Config{
				Field: "attributes",
			},
			expectedErr: "only unrolling logs from a body slice is currently supported",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
