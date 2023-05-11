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

// package removeemptyvalueprocessor provides a processor that removes empty values from telemetry data
package removeemptyvalueprocessor

// Config is the configuration for the processor
type Config struct {
	RemoveNulls              bool     `mapstructure:"remove_nulls"`
	RemoveEmptyLists         bool     `mapstructure:"remove_empty_lists"`
	RemoveEmptyMaps          bool     `mapstructure:"remove_empty_maps"`
	EnableResourceAttributes bool     `mapstructure:"enable_resource_attributes"`
	EnableAttributes         bool     `mapstructure:"enable_attributes"`
	EnableLogBody            bool     `mapstructure:"enable_log_body"`
	EmptyStringValues        []string `mapstructure:"empty_string_values"`
}

// Validate validates the processor configuration
func (cfg Config) Validate() error {
	return nil
}
