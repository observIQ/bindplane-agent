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

package logdeduplicationprocessor

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

const (
	// fieldDelimiter is the delimiter used to split a field key into its parts.
	fieldDelimiter = "."

	// fieldEscapeKeyReplacement is the string used to temporarily replace escaped delimters while splitting a field key.
	fieldEscapeKeyReplacement = "{TEMP_REPLACE}"
)

// fieldRemover handles removes match fields from log records
type fieldRemover struct {
	fields []*field
}

// field represents a field and it's compount key to match on
type field struct {
	keyParts []string
}

// newFieldRemover creates a new field remover based on the pased in field keys
func newFieldRemover(fieldKeys []string) *fieldRemover {
	fe := &fieldRemover{
		fields: make([]*field, 0, len(fieldKeys)),
	}

	for _, f := range fieldKeys {
		fe.fields = append(fe.fields, &field{
			keyParts: splitField(f),
		})
	}

	return fe
}

// RemoveFields removes any body or attribute fields that match in the log record
func (fe *fieldRemover) RemoveFields(logRecord plog.LogRecord) {
	for _, field := range fe.fields {
		field.removeField(logRecord)
	}
}

// removeField removes the field from the log record if it exists
func (f *field) removeField(logRecord plog.LogRecord) {
	// Get first key part
	firstPart, remainingParts := f.keyParts[0], f.keyParts[1:]

	switch firstPart {
	case bodyField:
		// If body is a map then recurse through to remove the field
		if logRecord.Body().Type() == pcommon.ValueTypeMap {
			removeFieldFromMap(logRecord.Body().Map(), remainingParts)
		}
	case attributeField:
		// Remove all attributes
		if len(remainingParts) == 0 {
			logRecord.Attributes().Clear()
			return
		}

		// Recurse through map and remove fields
		removeFieldFromMap(logRecord.Attributes(), remainingParts)
	}
}

// removeFieldFromMap recruses through the map and removes the field if it's found.
func removeFieldFromMap(valueMap pcommon.Map, keyParts []string) {
	// Get the next part of the key
	nextKeyPart, remainingParts := keyParts[0], keyParts[1:]

	// Look for the value associated with the next key part.
	// If we don't find it then return
	value, ok := valueMap.Get(nextKeyPart)
	if !ok {
		return
	}

	// No more key parts that means we have found the value and remove it
	if len(remainingParts) == 0 {
		valueMap.Remove(nextKeyPart)
		return
	}

	// If the value is a map then recurse through with the remaining parts
	if value.Type() == pcommon.ValueTypeMap {
		removeFieldFromMap(value.Map(), remainingParts)
	}
}

// splitField splits a field key into its parts.
// It replaces escaped delimiters with the full delimter after splitting.
func splitField(fieldKey string) []string {
	escapedKey := strings.ReplaceAll(fieldKey, fmt.Sprintf("\\%s", fieldDelimiter), fieldEscapeKeyReplacement)
	keyParts := strings.Split(escapedKey, fieldDelimiter)

	// Replace the temporarily escaped delimiters with the actual delimiter.
	for i := range keyParts {
		keyParts[i] = strings.ReplaceAll(keyParts[i], fieldEscapeKeyReplacement, fieldDelimiter)
	}

	return keyParts
}
