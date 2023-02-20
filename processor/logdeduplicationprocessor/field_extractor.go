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

// fieldExtractor handles extracting match fields from log records
type fieldExtractor struct {
	fields []*field
}

// field represents a field and it's compount key to match on
type field struct {
	keyParts []string
}

// newFieldExtractor creates a new field extractor based on the pased in field keys
func newFieldExtractor(fieldKeys []string) (*fieldExtractor, error) {
	fe := &fieldExtractor{
		fields: make([]*field, 0, len(fieldKeys)),
	}

	for _, f := range fieldKeys {
		fe.fields = append(fe.fields, &field{
			keyParts: splitField(f),
		})
	}

	return fe, nil
}

// HasFields returns true if all fields are present in the log record
func (fe *fieldExtractor) HasFields(logRecord plog.LogRecord) bool {
	for _, field := range fe.fields {
		if _, ok := field.extractField(logRecord); !ok {
			return false
		}
	}

	return true
}

// extractField extracts the field from the log record
func (f *field) extractField(logRecord plog.LogRecord) (pcommon.Value, bool) {
	// Get first key part
	firstPart, remainingParts := f.keyParts[0], f.keyParts[1:]

	switch firstPart {
	case bodyField:
		if len(remainingParts) == 0 {
			return logRecord.Body(), true
		} else if logRecord.Body().Type() != pcommon.ValueTypeMap {
			// Body is not a map and we have more keys to recurse through so return failure case
			return pcommon.NewValueEmpty(), false
		}

		// Recurse into the body
		return extractFieldFromMap(logRecord.Body().Map(), remainingParts)
	case attributeField:
		if len(remainingParts) == 0 {
			return pcommon.Value(logRecord.Attributes()), true
		}

		return extractFieldFromMap(logRecord.Attributes(), remainingParts)
	default:
		// Should not get here due to protections on config validation but just for completeness
		// TODO log
		return pcommon.NewValueEmpty(), false
	}

}

// extractFieldFromMap recruses through the map and extracts the field.
// If the key parts does not fully match a path then the function returns false.
func extractFieldFromMap(valueMap pcommon.Map, keyParts []string) (pcommon.Value, bool) {
	// Get the next part of the key
	nextKeyPart, remainingParts := keyParts[0], keyParts[1:]

	// Look for the value associated with the next key part.
	// If we don't find it return failure case.
	value, ok := valueMap.Get(nextKeyPart)
	if !ok {
		return pcommon.NewValueEmpty(), false
	}

	// No more key parts that means we have found the value
	if len(remainingParts) == 0 {
		return value, true
	}

	// We have more key parts to extract and the value is not a map type so return failure case
	if value.Type() != pcommon.ValueTypeMap {
		return pcommon.NewValueEmpty(), false
	}

	// Recurse into map with reminding key parts
	return extractFieldFromMap(value.Map(), remainingParts)
}

// splitField splits a field key into its parts.
// It replaces escaped delimiters with the full delimter after splitting.
func splitField(fieldKey string) []string {
	escapedKey := strings.ReplaceAll(fieldKey, fmt.Sprintf("\\%s", fieldDelimiter), fieldEscapeKeyReplacement)
	keyParts := strings.Split(escapedKey, fieldDelimiter)

	// Replace the temporarily escaped delimiters with the actual delimiter.
	for _, part := range keyParts {
		part = strings.ReplaceAll(part, fieldEscapeKeyReplacement, fieldDelimiter)
	}

	return keyParts
}
