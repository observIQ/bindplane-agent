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

// Package throughputwrapper represents a wrapper that wraps receivers and measures throughput
package throughputwrapper

import (
	"context"

	"go.opencensus.io/stats/view"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
)

// RegisterMetricViews unregisters old metric views if they exist and registers new ones
func RegisterMetricViews() error {
	views := metricViews()
	view.Unregister(views...)
	return view.Register(views...)
}

// WrapReceiverFactory creates a wrapper factory that around the passed in factory. Injecting consumers to measure output from the passed in receiver.
func WrapReceiverFactory(receiverFactory component.ReceiverFactory) component.ReceiverFactory {
	opts := make([]component.ReceiverFactoryOption, 0, 3)

	// Wrap the metric receiver creation func
	opts = append(opts, component.WithMetricsReceiver(
		wrapCreateMetricsReceiverFunc(receiverFactory.CreateMetricsReceiver), receiverFactory.MetricsReceiverStability()),
	)

	// Wrap the log receiver creation func
	opts = append(opts, component.WithLogsReceiver(
		wrapCreateLogReceiverFunc(receiverFactory.CreateLogsReceiver), receiverFactory.LogsReceiverStability()),
	)

	// Wrap the trace receiver creation func
	opts = append(opts, component.WithTracesReceiver(
		wrapCreateTraceReceiverFunc(receiverFactory.CreateTracesReceiver), receiverFactory.TracesReceiverStability()),
	)

	return component.NewReceiverFactory(
		receiverFactory.Type(),
		receiverFactory.CreateDefaultConfig,
		opts...,
	)
}

func wrapCreateMetricsReceiverFunc(createMetricsReceiverFunc component.CreateMetricsReceiverFunc) component.CreateMetricsReceiverFunc {
	return func(ctx context.Context,
		set component.ReceiverCreateSettings,
		rConf component.ReceiverConfig,
		nextConsumer consumer.Metrics,
	) (component.MetricsReceiver, error) {
		wrappedConsumer := newMetricConsumer(set.Logger, rConf.ID().String(), nextConsumer)
		return createMetricsReceiverFunc(ctx, set, rConf, wrappedConsumer)
	}
}

func wrapCreateLogReceiverFunc(createLogsReceiverFunc component.CreateLogsReceiverFunc) component.CreateLogsReceiverFunc {
	return func(ctx context.Context,
		set component.ReceiverCreateSettings,
		rConf component.ReceiverConfig,
		nextConsumer consumer.Logs,
	) (component.LogsReceiver, error) {
		wrappedConsumer := newLogConsumer(set.Logger, rConf.ID().String(), nextConsumer)
		return createLogsReceiverFunc(ctx, set, rConf, wrappedConsumer)
	}
}

func wrapCreateTraceReceiverFunc(createTracesReceiverFunc component.CreateTracesReceiverFunc) component.CreateTracesReceiverFunc {
	return func(ctx context.Context,
		set component.ReceiverCreateSettings,
		rConf component.ReceiverConfig,
		nextConsumer consumer.Traces,
	) (component.TracesReceiver, error) {
		wrappedConsumer := newTraceConsumer(set.Logger, rConf.ID().String(), nextConsumer)

		return createTracesReceiverFunc(ctx, set, rConf, wrappedConsumer)
	}
}
