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

package snapshot

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest/pmetrictest"
	"github.com/stretchr/testify/require"
)

var unfilteredMetricsDir = filepath.Join("testdata", "metrics", "before")
var filteredMetricsDir = filepath.Join("testdata", "metrics", "after")

func TestFilterMetrics(t *testing.T) {
	testCases := []struct {
		name            string
		fileIn          string
		query           *string
		minTimestamp    *time.Time
		expectedFileOut string
	}{
		{
			name:            "Matches attribute value (gauge)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "host-metrics.yaml"),
			query:           asPtr("cool-attribute-value"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-val-gauge.yaml"),
		},
		{
			name:            "Matches attribute key (gauge)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "host-metrics.yaml"),
			query:           asPtr("cool-attribute-key"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-key-gauge.yaml"),
		},
		{
			name:            "Matches attribute value (sum)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "host-metrics.yaml"),
			query:           asPtr("extra-sum-attr-value"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-val-sum.yaml"),
		},
		{
			name:            "Matches attribute key (sum)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "host-metrics.yaml"),
			query:           asPtr("extra-sum-attr-key"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-key-sum.yaml"),
		},
		{
			name:            "Matches resource attribute key",
			fileIn:          filepath.Join(unfilteredMetricsDir, "host-metrics.yaml"),
			query:           asPtr("extra-resource-attr-key"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-resource-attr-key.yaml"),
		},
		{
			name:            "Matches resource attribute value",
			fileIn:          filepath.Join(unfilteredMetricsDir, "host-metrics.yaml"),
			query:           asPtr("extra-resource-attr-value"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-resource-attr-val.yaml"),
		},
		{
			name:            "Matches metric name",
			fileIn:          filepath.Join(unfilteredMetricsDir, "host-metrics.yaml"),
			query:           asPtr("system.cpu.load_average.1m"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-metric-name.yaml"),
		},
		{
			name:   "Filters older than minimum timestamp",
			fileIn: filepath.Join(unfilteredMetricsDir, "host-metrics.yaml"),
			query:  asPtr("/dev"),
			// note when regenning goldens. golden.WriteMetrics
			// completely obliterates timestamps, so the regenerated golden will have the wrong timestamp.
			// You'll have to manually fix that.
			minTimestamp:    asPtr(time.Unix(0, 3000000)),
			expectedFileOut: filepath.Join(filteredMetricsDir, "filters-timestamp.yaml"),
		},
		{
			name:            "Matches attribute value (histogram)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "histogram.yaml"),
			query:           asPtr("dev-2"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-val-histogram.yaml"),
		},
		{
			name:            "Matches attribute key (histogram)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "histogram.yaml"),
			query:           asPtr("prod-machine"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-key-histogram.yaml"),
		},
		{
			name:            "Matches attribute value (exponential histogram)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "exp-histogram.yaml"),
			query:           asPtr("dev-2"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-val-exp-histogram.yaml"),
		},
		{
			name:            "Matches attribute key (exponential histogram)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "exp-histogram.yaml"),
			query:           asPtr("prod-machine"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-key-exp-histogram.yaml"),
		},
		{
			name:            "Matches attribute value (summary)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "summary.yaml"),
			query:           asPtr("dev-2"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-val-summary.yaml"),
		},
		{
			name:            "Matches attribute key (summary)",
			fileIn:          filepath.Join(unfilteredMetricsDir, "summary.yaml"),
			query:           asPtr("prod-machine"),
			expectedFileOut: filepath.Join(filteredMetricsDir, "matches-attr-key-summary.yaml"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metricsIn, err := golden.ReadMetrics(tc.fileIn)
			require.NoError(t, err)

			metricsOut := filterMetrics(metricsIn, tc.query, tc.minTimestamp)

			expectedMetricsOut, err := golden.ReadMetrics(tc.expectedFileOut)
			require.NoError(t, err)

			err = pmetrictest.CompareMetrics(expectedMetricsOut, metricsOut)
			require.NoError(t, err)
		})
	}
}
