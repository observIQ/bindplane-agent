package aggregationprocessor

import (
	"testing"
	"time"

	"github.com/observiq/observiq-otel-collector/processor/aggregationprocessor/internal/aggregate"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
)

func TestConfig_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		input       Config
		expectedErr string
	}{
		{
			name: "Config with all aggregations",
			input: Config{
				Interval: 5 * time.Second,
				Include:  "^.*$",
				Aggregations: []AggregateConfig{
					{
						Type: aggregate.AggregateTypeAvg,
					},
					{
						Type: aggregate.AggregateTypeMin,
					},
					{
						Type: aggregate.AggregateTypeMax,
					},
					{
						Type: aggregate.AggregateTypeLast,
					},
					{
						Type: aggregate.AggregateTypeFirst,
					},
				},
			},
		},
		{
			name: "Config with no aggregations",
			input: Config{
				Interval:     5 * time.Second,
				Include:      "^.*$",
				Aggregations: []AggregateConfig{},
			},
			expectedErr: "at least one aggregation must be specified",
		},
		{
			name: "Config with invalid regex",
			input: Config{
				Interval: 5 * time.Second,
				Include:  "^(",
				Aggregations: []AggregateConfig{
					{
						Type: aggregate.AggregateTypeAvg,
					},
				},
			},
			expectedErr: "`include` regex must be valid",
		},
		{
			name: "Config with invalid interval",
			input: Config{
				Interval: -5 * time.Second,
				Include:  "^.*$",
				Aggregations: []AggregateConfig{
					{
						Type: aggregate.AggregateTypeAvg,
					},
				},
			},
			expectedErr: "aggregation interval must be positive",
		},
		{
			name: "Config with invalid aggregation type",
			input: Config{
				Interval: 5 * time.Second,
				Include:  "^.*$",
				Aggregations: []AggregateConfig{
					{
						Type: aggregate.AggregateType("invalid"),
					},
				},
			},
			expectedErr: "invalid aggregate type for `type`: invalid",
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

func TestAggregateConfig_MetricNameString(t *testing.T) {
	t.Run("metric name is not specified", func(t *testing.T) {
		metricName := AggregateConfig{
			Type:       aggregate.AggregateTypeAvg,
			MetricName: "",
		}.MetricNameString()
		require.Equal(t, "$0", metricName)
	})

	t.Run("metric name is specified", func(t *testing.T) {
		metricName := AggregateConfig{
			Type:       aggregate.AggregateTypeAvg,
			MetricName: "test.metric",
		}.MetricNameString()
		require.Equal(t, "test.metric", metricName)
	})
}

func TestValidStruct(t *testing.T) {
	componenttest.CheckConfigStruct(&Config{})
}
