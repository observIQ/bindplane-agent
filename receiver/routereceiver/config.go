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

// Package routereceiver provides a receiver that receives telemetry from other components.
package routereceiver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

// Config is the config of the route receiver.
type Config struct {
	config.ReceiverSettings `mapstructure:",squash"`
}

// createDefaultConfig returns the default config for the route receiver.
func createDefaultConfig() component.Config {
	return &Config{
		ReceiverSettings: config.NewReceiverSettings(component.NewID(typeStr)),
	}
}
