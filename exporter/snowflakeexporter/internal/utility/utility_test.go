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

package utility

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.uber.org/zap"
)

func TestConvertAttributesToString(t *testing.T) {
	testCases := []struct {
		desc       string
		attributes map[string]any
		expected   string
	}{
		{
			desc: "simple",
			attributes: map[string]any{
				"k1": "v1",
				"k2": "v2",
			},
			expected: "{\"k1\":\"v1\",\"k2\":\"v2\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			m := pcommon.NewMap()
			require.NoError(t, m.FromRaw(tc.attributes))
			attributeString := ConvertAttributesToString(m, zap.NewNop())
			require.Equal(t, tc.expected, attributeString)
		})
	}
}
