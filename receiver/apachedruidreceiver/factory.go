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

package apachedruidreceiver

import (
	"context"

	"github.com/observiq/bindplane-agent/receiver/apachedruidreceiver/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	typeStr   = "apachedruid"
	stability = component.StabilityLevelDevelopment
)

// NewFactory creates a factory for an Apache Druid receiver
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, stability),
	)
}

// createMetricsReceiver creates an Apache Druid receiver with a metrics consumer
func createMetricsReceiver(
	_ context.Context,
	settings receiver.CreateSettings,
	config component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	cfg := config.(*Config)
	return newMetricsReceiver(settings, cfg, consumer)
}

// createDefaultConfig creates a default config for an Apache Druid receiver
func createDefaultConfig() component.Config {
	return &Config{
		Metrics: MetricsConfig{
			MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
		},
	}
}
