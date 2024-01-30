// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snapshot

import (
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

func queryMatchesValue(v pcommon.Value, query string) bool {
	switch v.Type() {
	case pcommon.ValueTypeMap:
		// Recursively query for a match in the map (depth first search)
		return queryMatchesMap(v.Map(), query)
	case pcommon.ValueTypeSlice:
		// Iterate each element of the slice
		return queryMatchesSlice(v.Slice(), query)
	case pcommon.ValueTypeEmpty:
		// Cannot match empty value
		return false
	default:
		// We might be able to actually get away with just doing this for slices/maps, but could lead to
		// weird edgecases since those slices/maps would be json-ified
		// Note: Bytes will be base64 encoded and searched that way.
		return strings.Contains(v.AsString(), query)
	}
}

// func matchesAttributes(l plog.LogRecord, search string) {}
func queryMatchesMap(m pcommon.Map, query string) bool {
	matches := false

	m.Range(func(k string, v pcommon.Value) bool {
		// check if key matches
		matches = strings.Contains(k, query)
		if matches {
			// Return false to cancel iterating, since we know this map matches
			return false
		}

		// Check if the value matches
		matches = queryMatchesValue(v, query)
		if matches {
			// Return false to cancel iterating, since we know this map matches
			return false
		}

		// Continue iterating since we haven't found a match
		return true
	})

	return matches
}

func queryMatchesSlice(s pcommon.Slice, query string) bool {
	for i := 0; i < s.Len(); i++ {
		elem := s.At(i)

		if queryMatchesValue(elem, query) {
			return true
		}
	}

	return false
}
