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

package removeemptyvaluesprocessor

import (
	"context"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type emptyValueProcessor struct {
	logger               *zap.Logger
	c                    Config
	excludeResourceKeys  []MapKey
	excludeAttributeKeys []MapKey
	excludeBodyKeys      []MapKey
}

func newEmptyValueProcessor(logger *zap.Logger, cfg Config) *emptyValueProcessor {
	var excludeResourceKeys []MapKey
	var excludeAttributeKeys []MapKey
	var excludeBodyKeys []MapKey

	for _, mapKey := range cfg.ExcludeKeys {
		switch mapKey.Field {
		case AttributesField:
			excludeAttributeKeys = append(excludeAttributeKeys, mapKey)
		case ResourceField:
			excludeResourceKeys = append(excludeResourceKeys, mapKey)
		case BodyField:
			excludeBodyKeys = append(excludeBodyKeys, mapKey)
		}
	}

	return &emptyValueProcessor{
		logger:               logger,
		c:                    cfg,
		excludeResourceKeys:  excludeResourceKeys,
		excludeAttributeKeys: excludeAttributeKeys,
		excludeBodyKeys:      excludeBodyKeys,
	}
}

func (evp *emptyValueProcessor) processTraces(_ context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	resourceSpans := td.ResourceSpans()
	for i := 0; i < resourceSpans.Len(); i++ {
		resourceSpan := resourceSpans.At(i)
		scopeSpans := resourceSpan.ScopeSpans()

		if evp.c.EnableResourceAttributes {
			cleanMap(resourceSpan.Resource().Attributes(), evp.c, evp.excludeResourceKeys)
		}

		if !evp.c.EnableAttributes {
			// Skip loops for attributes if we don't need to clean them.
			continue
		}

		for j := 0; j < scopeSpans.Len(); j++ {
			scopeSpan := scopeSpans.At(j)
			spans := scopeSpan.Spans()

			for k := 0; k < spans.Len(); k++ {
				span := spans.At(k)
				cleanMap(span.Attributes(), evp.c, evp.excludeAttributeKeys)
			}
		}
	}

	return td, nil
}

func (evp *emptyValueProcessor) processLogs(_ context.Context, ld plog.Logs) (plog.Logs, error) {
	resourceLogs := ld.ResourceLogs()
	for i := 0; i < resourceLogs.Len(); i++ {
		resourceLog := resourceLogs.At(i)
		scopeLogs := resourceLog.ScopeLogs()

		if evp.c.EnableResourceAttributes {
			cleanMap(resourceLog.Resource().Attributes(), evp.c, evp.excludeResourceKeys)
		}

		for j := 0; j < scopeLogs.Len(); j++ {
			scopeLog := scopeLogs.At(j)
			logRecords := scopeLog.LogRecords()

			for k := 0; k < logRecords.Len(); k++ {
				logRecord := logRecords.At(k)
				if evp.c.EnableAttributes {
					cleanMap(logRecord.Attributes(), evp.c, evp.excludeAttributeKeys)
				}

				if evp.c.EnableLogBody {
					cleanLogBody(logRecord, evp.c, evp.excludeBodyKeys)
				}
			}
		}
	}

	return ld, nil
}

func (evp *emptyValueProcessor) processMetrics(_ context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	resourceMetrics := md.ResourceMetrics()
	for i := 0; i < resourceMetrics.Len(); i++ {
		resourceMetric := resourceMetrics.At(i)
		scopeMetrics := resourceMetric.ScopeMetrics()

		if evp.c.EnableResourceAttributes {
			cleanMap(resourceMetric.Resource().Attributes(), evp.c, evp.excludeResourceKeys)
		}

		if !evp.c.EnableAttributes {
			// Skip loops for attributes if we don't need to clean them.
			continue
		}

		for j := 0; j < scopeMetrics.Len(); j++ {
			scopeMetric := scopeMetrics.At(j)
			metrics := scopeMetric.Metrics()

			for k := 0; k < metrics.Len(); k++ {
				metric := metrics.At(k)
				cleanMetricAttrs(metric, evp.c, evp.excludeAttributeKeys)
			}
		}
	}
	return md, nil
}

// cleanMap removes empty values from the map, as defined by the config.
func cleanMap(m pcommon.Map, c Config, excludeKeys []MapKey) {
	m.RemoveIf(func(s string, v pcommon.Value) bool {
		for _, mk := range excludeKeys {
			if mk.Key == s {
				return false
			}
		}

		switch v.Type() {
		case pcommon.ValueTypeEmpty:
			return c.RemoveNulls
		case pcommon.ValueTypeMap:
			subMap := v.Map()
			cleanMap(subMap, c, trimMapKeyPrefix(s, excludeKeys))
			return c.RemoveEmptyMaps && subMap.Len() == 0
		case pcommon.ValueTypeSlice:
			return c.RemoveEmptyLists && v.Slice().Len() == 0
		case pcommon.ValueTypeStr:
			return shouldFilterString(v.Str(), c.EmptyStringValues)
		}

		return false
	})
}

// trimMapKeyPrefix returns the provided keys with the specified prefix removed.
// Any keys that don't have the prefix are removed from the returned list.
func trimMapKeyPrefix(prefix string, keys []MapKey) []MapKey {
	outKeys := make([]MapKey, 0, len(keys))
	for _, mk := range keys {
		trimmedKey, found := strings.CutPrefix(mk.Key, prefix+".")
		if !found {
			// prefix was not found, so this key does not belong to the submap.
			continue
		}

		outKeys = append(outKeys, MapKey{
			Field: mk.Field,
			Key:   trimmedKey,
		})
	}

	return outKeys
}

// shouldFilterString returns true if the given string should be considered an "empty" value,
// according to the config.
func shouldFilterString(s string, emptyValues []string) bool {
	for _, filteredString := range emptyValues {
		if s == filteredString {
			return true
		}
	}

	return false
}

// cleanMetricAttrs removes any attributes that should be considered empty from all the datapoints in the metrics.
func cleanMetricAttrs(metric pmetric.Metric, c Config, keys []MapKey) {
	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		dps := metric.Gauge().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c, keys)
		}

	case pmetric.MetricTypeHistogram:
		dps := metric.Histogram().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c, keys)
		}
	case pmetric.MetricTypeSum:
		dps := metric.Sum().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c, keys)
		}
	case pmetric.MetricTypeSummary:
		dps := metric.Summary().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c, keys)
		}
	case pmetric.MetricTypeExponentialHistogram:
		dps := metric.ExponentialHistogram().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c, keys)
		}
	default:
		// skip metric if None or unknown type
	}
}

// cleanLogBody removes empty values from the log body.
func cleanLogBody(lr plog.LogRecord, c Config, keys []MapKey) {
	body := lr.Body()
	switch body.Type() {
	case pcommon.ValueTypeMap:
		bodyMap := body.Map()
		cleanMap(bodyMap, c, keys)
		if c.RemoveEmptyMaps && bodyMap.Len() == 0 {
			pcommon.NewValueEmpty().CopyTo(body)
		}
	case pcommon.ValueTypeSlice:
		bodySlice := body.Slice()
		if c.RemoveEmptyLists && bodySlice.Len() == 0 {
			pcommon.NewValueEmpty().CopyTo(body)
		}
	case pcommon.ValueTypeStr:
		if shouldFilterString(body.Str(), c.EmptyStringValues) {
			pcommon.NewValueEmpty().CopyTo(body)
		}
	}
}
