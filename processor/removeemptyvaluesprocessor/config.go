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

// Package removeemptyvaluesprocessor provides a processor that removes empty values from telemetry data
package removeemptyvaluesprocessor

import (
	"bytes"
	"fmt"
)

// valid fields that can be referenced in the MapKey's field
const (
	attributesField = "attributes"
	resourceField   = "resource"
	bodyField       = "body"
)

// MapKey represents a key into a particular map (denoted by field)
type MapKey struct {
	field string
	key   string
}

// UnmarshalText unmarshals the given []byte into a MapKey.
// The format of the key is "<field>.<path-to-key>"
func (m *MapKey) UnmarshalText(text []byte) error {
	field, key, found := bytes.Cut(text, []byte("."))
	if !found {
		return fmt.Errorf("failed to determine field: %s", text)
	}

	if len(key) == 0 {
		return fmt.Errorf("key part of (%s) must be non-zero in length", text)
	}

	for _, validField := range []string{attributesField, resourceField, bodyField} {
		if validField == string(field) {
			// this key indexes into a valid field, and therefore
			m.key = string(key)
			m.field = string(field)
			return nil
		}
	}

	return fmt.Errorf("invalid field (%s), must be one of attributes, resource, or body", field)
}

// Config is the configuration for the processor
type Config struct {
	RemoveNulls              bool     `mapstructure:"remove_nulls"`
	RemoveEmptyLists         bool     `mapstructure:"remove_empty_lists"`
	RemoveEmptyMaps          bool     `mapstructure:"remove_empty_maps"`
	EnableResourceAttributes bool     `mapstructure:"enable_resource_attributes"`
	EnableAttributes         bool     `mapstructure:"enable_attributes"`
	EnableLogBody            bool     `mapstructure:"enable_log_body"`
	EmptyStringValues        []string `mapstructure:"empty_string_values"`
	ExcludeKeys              []MapKey `mapstructure:"exclude_keys"`
}

// Validate validates the processor configuration
func (cfg Config) Validate() error {
	return nil
}
