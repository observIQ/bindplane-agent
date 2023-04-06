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

package spancountprocessor

import (
	"context"
	"fmt"

	"github.com/observiq/observiq-otel-collector/expr"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

const (
	// typeStr is the value of the "type" key in configuration.
	typeStr = "spancount"

	// stability is the current state of the processor.
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a new factory for the processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		typeStr,
		createDefaultConfig,
		processor.WithTraces(createTracesProcessor, stability),
	)
}

// createTracesProcessor creates a log processor.
func createTracesProcessor(_ context.Context, params processor.CreateSettings, cfg component.Config, consumer consumer.Traces) (processor.Traces, error) {
	processorCfg, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type: %+v", cfg)
	}

	match, err := expr.CreateBoolExpression(processorCfg.Match)
	if err != nil {
		return nil, fmt.Errorf("invalid match expression: %w", err)
	}

	attrs, err := expr.CreateExpressionMap(processorCfg.Attributes)
	if err != nil {
		return nil, fmt.Errorf("invalid attribute expression: %w", err)
	}

	return newProcessor(processorCfg, consumer, match, attrs, params.Logger), nil
}
