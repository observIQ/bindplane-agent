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

package googlemanagedprometheusexporter

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/processor"
)

// googleManagedPrometheusExporter is a google managed prometheus exporter wrapped with additional functionality
type googleManagedPrometheusExporter struct {
	metricsProcessors []processor.Metrics
	metricsExporter   exporter.Metrics
	metricsConsumer   consumer.Metrics
}

// ConsumeMetrics consumes metrics
func (e *googleManagedPrometheusExporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	if e.metricsConsumer == nil {
		return nil
	}

	return e.metricsConsumer.ConsumeMetrics(ctx, md)
}

// Capabilities returns the capabilities of the exporter
func (e *googleManagedPrometheusExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Start starts the exporter
func (e *googleManagedPrometheusExporter) Start(ctx context.Context, host component.Host) error {
	if e.metricsExporter != nil {
		if err := e.metricsExporter.Start(ctx, host); err != nil {
			return fmt.Errorf("failed to start metrics exporter: %w", err)
		}
	}

	for _, processor := range e.metricsProcessors {
		if err := processor.Start(ctx, host); err != nil {
			return fmt.Errorf("failed to start metrics processor: %w", err)
		}
	}

	return nil
}

// Shutdown will shutdown the exporter
func (e *googleManagedPrometheusExporter) Shutdown(ctx context.Context) error {
	for i := len(e.metricsProcessors) - 1; i >= 0; i-- {
		if err := e.metricsProcessors[i].Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown metrics processor: %w", err)
		}
	}

	if e.metricsExporter != nil {
		if err := e.metricsExporter.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown metrics exporter: %w", err)
		}
	}

	return nil
}
