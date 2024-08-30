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

// Package marshalprocessor provides a processor that marshals logs to a specified format.
package marshalprocessor

import (
	"errors"
	"strings"
)

var errInvalidMarshalTo = errors.New("marshal_to must be JSON, XML, or KV")

// Config is the configuration for the processor
type Config struct {
	// MarshalTo is either JSON, XML, or KV
	MarshalTo string `mapstructure:"marshal_to"`
}

// Validate validates the processor configuration
func (cfg Config) Validate() error {
	// Validate MarshalTo choice
	switch strings.ToUpper(cfg.MarshalTo) {
	case "JSON":
		return nil
	case "XML":
		return nil
	case "KV":
		return nil
	default:
		return errInvalidMarshalTo
	}
}
