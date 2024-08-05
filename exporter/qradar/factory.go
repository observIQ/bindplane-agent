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

package qradar

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/pdata/plog"
)

var otlpFactory = otlphttpexporter.NewFactory()

// componentType is the type of the qradar exporter
var componentType = component.MustNewType("qradar")

const (
	// The stability level of the exporter. Matches the current exporter in contrib
	stability = component.StabilityLevelBeta
)

// NewFactory creates a factory for the qradar exporter
func NewFactory(collectorVersion string) exporter.Factory {
	return exporter.NewFactory(
		componentType,
		createDefaultConfig(collectorVersion),
		exporter.WithLogs(createLogsExporter, stability),
	)
}

type myExporter struct {
	httpExporter exporter.Logs
}

func (e *myExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {

	// distribute resource across records

	resourceLogs := ld.ResourceLogs()
	for i := 0; i < resourceLogs.Len(); i++ {
		rl := resourceLogs.At(i)
		rl.Resource()
	}

	// send to http exporter
	e.httpExporter.ConsumeLogs(ctx, ld)
	return nil
}

func createLogsExporter(ctx context.Context, settings exporter.Settings, cfg component.Config) (exporter.Logs, error) {
	exporterConfig := cfg.(*Config)

	customhttp, err := otlpFactory.CreateLogsExporter(ctx, settings, exporterConfig.OTLPConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs exporter: %w", err)
	}

	return customhttp, nil
}
