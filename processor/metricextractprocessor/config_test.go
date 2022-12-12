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

package metricextractprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
)

func TestCreateDefaultProcessorConfig(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	require.Equal(t, defaultMatch, cfg.Match)
	require.Equal(t, defaultMetricName, cfg.MetricName)
	require.Equal(t, defaultMetricUnit, cfg.MetricUnit)
	require.Equal(t, component.NewID(typeStr), cfg.ProcessorSettings.ID())
}

func TestConfigValidate(t *testing.T) {
	var testCases = []struct {
		name     string
		config   *Config
		expected error
	}{
		{
			name: "valid config",
			config: &Config{
				Match:      "true",
				Extract:    "message",
				MetricName: "metric",
				MetricUnit: "unit",
				MetricType: gaugeDoubleType,
			},
			expected: nil,
		},
		{
			name: "invalid metric type",
			config: &Config{
				Match:      "true",
				Extract:    "message",
				MetricName: "metric",
				MetricUnit: "unit",
				MetricType: "invalid",
			},
			expected: errMetricTypeInvalid,
		},
		{
			name: "missing extract",
			config: &Config{
				Match:      "true",
				MetricName: "metric",
				MetricUnit: "unit",
				MetricType: gaugeDoubleType,
			},
			expected: errExtractMissing,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			require.Equal(t, tc.expected, err)
		})
	}
}
