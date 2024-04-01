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

package routereceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

// componentType is the value of the "type" key in configuration.
var componentType = component.MustNewType("route")

const (
	// stability is the current state of the receiver.
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a new factory for the receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		componentType,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, stability),
		receiver.WithLogs(createLogsReceiver, stability),
		receiver.WithTraces(createTracesReceiver, stability),
	)
}

// createMetricsReceiver creates a metric receiver.
func createMetricsReceiver(_ context.Context, set receiver.CreateSettings, _ component.Config, consumer consumer.Metrics) (receiver.Metrics, error) {
	receiver := createOrGetRoute(set.ID.Name())
	receiver.registerMetricConsumer(consumer)
	return receiver, nil
}

// createLogsReceiver creates a log receiver.
func createLogsReceiver(_ context.Context, set receiver.CreateSettings, _ component.Config, consumer consumer.Logs) (receiver.Logs, error) {
	receiver := createOrGetRoute(set.ID.Name())
	receiver.registerLogConsumer(consumer)
	return receiver, nil
}

// createTracesReceiver creates a trace receiver.
func createTracesReceiver(_ context.Context, set receiver.CreateSettings, _ component.Config, consumer consumer.Traces) (receiver.Traces, error) {
	receiver := createOrGetRoute(set.ID.Name())
	receiver.registerTraceConsumer(consumer)
	return receiver, nil
}
