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

// Package utility provides utility functions for the snowflakeexporter package to consolidate code
package utility

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/lib/pq"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// CreateDSN creates a DSN for connecting to Snowflake with the given config
// TODO add functionality for additional query params
func CreateDSN(username, password, accountID, database string) string {
	usernameEsc := url.QueryEscape(username)
	passwordEsc := url.QueryEscape(password)
	accountIDEsc := url.QueryEscape(accountID)
	databaseEsc := url.QueryEscape(database)
	return fmt.Sprintf(`%s:%s@%s/"%s"`, usernameEsc, passwordEsc, accountIDEsc, databaseEsc)
}

// ConvertAttributesToString converts the pcommon.Map into a JSON string representation
// this is due to a bug/lacking feature with the snowflake driver that prevents maps from being inserted into VARIANT & OBJECT columns
// github issue: https://github.com/snowflakedb/gosnowflake/issues/217
func ConvertAttributesToString(attributes pcommon.Map, logger *zap.Logger) string {
	m := make(map[string]string, attributes.Len())
	attributes.Range(func(k string, v pcommon.Value) bool {
		m[k] = v.AsString()
		return true
	})

	j, err := json.Marshal(m)
	if err != nil {
		logger.Warn("failed to marshal attribute map", zap.Error(err), zap.Any("map", m))
		return ""
	}
	return string(j)
}

// TraceIDToHexOrEmptyString returns a string representation of the TraceID in hex
func TraceIDToHexOrEmptyString(id pcommon.TraceID) string {
	if id.IsEmpty() {
		return ""
	}
	return hex.EncodeToString(id[:])
}

// SpanIDToHexOrEmptyString returns a string representation of the SpanID in hex
func SpanIDToHexOrEmptyString(id pcommon.SpanID) string {
	if id.IsEmpty() {
		return ""
	}
	return hex.EncodeToString(id[:])
}

// FlattenExemplars will flatten the given exemplars into slices of the individual fields
func FlattenExemplars(exemplars pmetric.ExemplarSlice, l *zap.Logger) (pq.StringArray, pq.StringArray, pq.StringArray, pq.StringArray, pq.Float64Array) {
	attributes := pq.StringArray{}
	timestamps := pq.StringArray{}
	traceIDs := pq.StringArray{}
	spanIDs := pq.StringArray{}
	values := pq.Float64Array{}

	for i := 0; i < exemplars.Len(); i++ {
		e := exemplars.At(i)
		attributes = append(attributes, ConvertAttributesToString(e.FilteredAttributes(), l))
		timestamps = append(timestamps, e.Timestamp().String())
		traceIDs = append(traceIDs, TraceIDToHexOrEmptyString(e.TraceID()))
		spanIDs = append(spanIDs, SpanIDToHexOrEmptyString(e.SpanID()))

		// convert Int value to Float64 so that values can be combined and still inserted as an array
		if e.ValueType() == pmetric.ExemplarValueTypeInt {
			values = append(values, float64(e.IntValue()))
		} else {
			values = append(values, e.DoubleValue())
		}
	}

	return attributes, timestamps, traceIDs, spanIDs, values
}
