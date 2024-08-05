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
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
)

type Config struct {
	// OTLP exporter configuration
	OTLPConfig *otlphttpexporter.Config `mapstructure:",squash"`
}

func createDefaultConfig(collectorVersion string) func() component.Config {
	return func() component.Config {
		return &Config{
			OTLPConfig: createDefaultOTLPConfig(collectorVersion),
		}
	}
}

// createDefaultGCPConfig creates a default GCP config
func createDefaultOTLPConfig(collectorVersion string) *otlphttpexporter.Config {
	factory := otlphttpexporter.NewFactory()
	config := factory.CreateDefaultConfig().(*otlphttpexporter.Config)

	return config
}
