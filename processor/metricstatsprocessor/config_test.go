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

package metricstatsprocessor

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/observiq/bindplane-otel-collector/processor/metricstatsprocessor/internal/stats"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	testCases := []struct {
		id       component.ID
		expected component.Config
	}{
		{
			id:       component.NewIDWithName(componentType, "defaults"),
			expected: createDefaultConfig(),
		},
		{
			id: component.NewIDWithName(componentType, ""),
			expected: &Config{
				Interval: 3 * time.Minute,
				Include:  `^test\.thing$$`,
				Stats: []stats.StatType{
					stats.LastType,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tc.id.String())
			require.NoError(t, err)
			require.NoError(t, sub.Unmarshal(cfg))

			assert.NoError(t, component.ValidateConfig(cfg))
			require.Equal(t, tc.expected, cfg)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		input       Config
		expectedErr string
	}{
		{
			name: "Config with all stat types",
			input: Config{
				Interval: 5 * time.Second,
				Include:  "^.*$",
				Stats: []stats.StatType{
					stats.AvgType,
					stats.MinType,
					stats.MaxType,
					stats.LastType,
					stats.FirstType,
				},
			},
		},
		{
			name: "Config with no stat types",
			input: Config{
				Interval: 5 * time.Second,
				Include:  "^.*$",
				Stats:    []stats.StatType{},
			},
			expectedErr: "at least one statistic must be specified in `stats`",
		},
		{
			name: "Config with default stat types",
			input: Config{
				Interval: 5 * time.Second,
				Include:  "^.*$",
			},
		},
		{
			name: "Config with invalid regex",
			input: Config{
				Interval: 5 * time.Second,
				Include:  "^(",
				Stats: []stats.StatType{
					stats.AvgType,
				},
			},
			expectedErr: "`include` regex must be valid",
		},
		{
			name: "Config with invalid interval",
			input: Config{
				Interval: -5 * time.Second,
				Include:  "^.*$",
				Stats: []stats.StatType{
					stats.AvgType,
				},
			},
			expectedErr: "interval must be positive",
		},
		{
			name: "Config with invalid stat type",
			input: Config{
				Interval: 5 * time.Second,
				Include:  "^.*$",
				Stats: []stats.StatType{
					stats.StatType("invalid"),
				},
			},
			expectedErr: "invalid statistic type for `type`: invalid",
		},
		{
			name: "Config with duplicate stat types",
			input: Config{
				Interval: 5 * time.Second,
				Include:  "^.*$",
				Stats: []stats.StatType{
					stats.AvgType,
					stats.AvgType,
				},
			},
			expectedErr: "each statistic type can only be specified once (avg specified more than once)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate()
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidStruct(t *testing.T) {
	require.NoError(t, componenttest.CheckConfigStruct(&Config{}))
}
