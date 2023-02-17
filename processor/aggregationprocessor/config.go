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

// Package aggregationprocessor provides a processor that samples pdata base level objects.
package aggregationprocessor

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/antonmedv/expr"
	col_expr "github.com/observiq/observiq-otel-collector/expr"
	"github.com/observiq/observiq-otel-collector/processor/aggregationprocessor/internal/aggregate"
)

// Config is the configuration for the processor
type Config struct {
	Interval time.Duration `mapstructure:"interval"`
	// Include is a regex that must match the metric name for it to be sampled.
	// Otherwise, the metric is passed through.
	Include string `mapstructure:"include"`
	// List of aggregations for the metric
	Aggregations []AggregateConfig `mapstructure:"aggregations"`
}

// AggregateConfig is a config that specifies which aggregations to perform for each incoming metric
type AggregateConfig struct {
	// Type of aggregation
	Type aggregate.AggregationType `mapstructure:"type"`
	// MetricNameExpr is an expression that gives a name for the metric. Defaults to `metric_name` (this is the original metric's name)
	MetricNameExprStr string `mapstructure:"metric_name_expression"`
	// Cached metric name expression
	metricNameExpr *col_expr.Expression
}

// Validate validate the config, returning an error explaining why it isn't if the config is invalid.
func (a AggregateConfig) Validate() error {
	if !a.Type.Valid() {
		return fmt.Errorf("invalid aggregate type for `type`: %s", a.Type)
	}

	_, err := a.MetricNameExpression()
	if err != nil {
		return fmt.Errorf("failed to parse metric_name_expression: %w", err)
	}

	return nil
}

const metricNameKey = "metric_name"

// MetricNameExpression returns a compiled expression for the given input
func (a AggregateConfig) MetricNameExpression() (*col_expr.Expression, error) {
	if a.metricNameExpr != nil {
		return a.metricNameExpr, nil
	}

	opts := []expr.Option{
		expr.AsKind(reflect.String),
		expr.Optimize(true),
		expr.Env(map[string]any{
			metricNameKey: "",
		}),
	}
	if a.MetricNameExprStr != "" {
		return col_expr.CreateExpression(a.MetricNameExprStr, opts...)
	}

	return col_expr.CreateExpression(metricNameKey, opts...)
}

// Validate validates the processor configuration
func (cfg Config) Validate() error {
	if _, err := regexp.Compile(cfg.Include); err != nil {
		return fmt.Errorf("`include` regex must be valid: %w", err)
	}

	if cfg.Interval <= 0 {
		return errors.New("aggregation interval must be positive")
	}

	// don't check aggregations if using defaults
	if cfg.Aggregations == nil {
		return nil
	}

	if len(cfg.Aggregations) == 0 {
		return errors.New("at least one aggregation must be specified")
	}

	for _, a := range cfg.Aggregations {
		if err := a.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// AggregationConfigs gets the default aggregation configs if none were specified, otherwise the specified aggregation configs
func (cfg Config) AggregationConfigs() []AggregateConfig {
	if cfg.Aggregations == nil {
		// fallback to defaults
		return []AggregateConfig{
			{
				Type:              aggregate.MinType,
				MetricNameExprStr: `metricNameKey + ".min"`,
			},
			{
				Type:              aggregate.MaxType,
				MetricNameExprStr: `metricNameKey + ".max"`,
			},
			{
				Type:              aggregate.AvgType,
				MetricNameExprStr: `metricNameKey + ".avg"`,
			},
		}
	}

	return cfg.Aggregations
}
