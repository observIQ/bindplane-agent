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

package serializeprocessor

import (
	"errors"
	"strings"
)

var errInvalidSerializeTo = errors.New("serialize_to must be JSON, XML, or KV")

// Config is the configuration for the processor
type Config struct {
	// SerializeTo is either JSON, XML, or KV
	SerializeTo	string `mapstructure:"serialize_to"`
}

// Validate validates the processor configuration
func (cfg Config) Validate() error {
	// Validate SerializeTo choice
	switch strings.ToUpper(cfg.SerializeTo) {
	case "JSON":
		return nil
	case "XML":
		return nil
	case "KV":
		return nil
	default:
		return errInvalidSerializeTo
	}
}
