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

package logcountprocessor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
)

const (
	// typeStr is the value of the "type" key in configuration.
	typeStr = "logcount"

	// stability is the current state of the receiver and processor.
	stability = component.StabilityLevelStable
)

// NewProcessorFactory creates a new factory for the processor.
func NewProcessorFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(
		typeStr,
		createDefaultProcessorConfig,
		component.WithLogsProcessor(createLogsProcessor, stability),
	)
}

// createLogsProcessor creates a log processor.
func createLogsProcessor(_ context.Context, params component.ProcessorCreateSettings, cfg component.ProcessorConfig, consumer consumer.Logs) (component.LogsProcessor, error) {
	processorCfg, ok := cfg.(*ProcessorConfig)
	if !ok {
		return nil, fmt.Errorf("invalid config type: %+v", cfg)
	}

	matchExpr, err := processorCfg.createMatchExpr()
	if err != nil {
		return nil, err
	}

	attrExprs, err := processorCfg.createAttrExprs()
	if err != nil {
		return nil, err
	}

	return newProcessor(processorCfg, consumer, matchExpr, attrExprs, params.Logger), nil
}

// NewReceiverFactory creates a new factory for the receiver.
func NewReceiverFactory() component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		createDefaultReceiverConfig,
		component.WithMetricsReceiver(createMetricsReceiver, stability),
	)
}

// createMetricsReceiver creates a metric receiver.
func createMetricsReceiver(_ context.Context, _ component.ReceiverCreateSettings, cfg component.ReceiverConfig, consumer consumer.Metrics) (component.MetricsReceiver, error) {
	return newReceiver(cfg.ID(), consumer), nil
}
