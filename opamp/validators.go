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

package opamp

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// ValidatorFunc function that takes in a config contents and validates
// Returns true if valid
type ValidatorFunc func([]byte) bool

// NoopValidator is a validator that always passes validation.
// Useful when unsure of what validator to use.
func NoopValidator(_ []byte) bool {
	return true
}

// NewYamlValidator creates a new Validator that checks does a yaml unmarshal against the target interface{}
func NewYamlValidator(target interface{}) ValidatorFunc {
	return func(b []byte) bool {
		if err := yaml.Unmarshal(b, target); err != nil {
			return false
		}

		return true
	}
}

// NewJSONValidator creates a new Validator that checks does a json unmarshal against the target interface{}
func NewJSONValidator(target interface{}) ValidatorFunc {
	return func(b []byte) bool {
		if err := json.Unmarshal(b, &target); err != nil {
			return false
		}

		return true
	}
}
