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

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

type emptyValueProcessor struct {
	logger *zap.Logger
	c      Config
}

func newEmptyValueProcessor(logger *zap.Logger, cfg Config) *emptyValueProcessor {
	return &emptyValueProcessor{
		logger: logger,
		c:      cfg,
	}
}

func (evp *emptyValueProcessor) processTraces(_ context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	resourceSpans := td.ResourceSpans()
	for i := 0; i < resourceSpans.Len(); i++ {
		resourceSpan := resourceSpans.At(i)
		scopeSpans := resourceSpan.ScopeSpans()

		if evp.c.EnableResourceAttributes {
			cleanMap(resourceSpan.Resource().Attributes(), evp.c)
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
				cleanMap(span.Attributes(), evp.c)
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
			cleanMap(resourceLog.Resource().Attributes(), evp.c)
		}

		for j := 0; j < scopeLogs.Len(); j++ {
			scopeLog := scopeLogs.At(j)
			logRecords := scopeLog.LogRecords()

			for k := 0; k < logRecords.Len(); k++ {
				logRecord := logRecords.At(k)
				if evp.c.EnableAttributes {
					cleanMap(logRecord.Attributes(), evp.c)
				}

				if evp.c.EnableLogBody {
					cleanLogBody(logRecord, evp.c)
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
			cleanMap(resourceMetric.Resource().Attributes(), evp.c)
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
				cleanMetricAttrs(metric, evp.c)
			}
		}
	}
	return md, nil
}

// cleanMap removes empty values from the map, as defined by the config.
func cleanMap(m pcommon.Map, c Config) {
	m.RemoveIf(func(s string, v pcommon.Value) bool {
		switch v.Type() {
		case pcommon.ValueTypeEmpty:
			return c.RemoveNulls
		case pcommon.ValueTypeMap:
			subMap := v.Map()
			cleanMap(subMap, c)
			return c.RemoveEmptyMaps && subMap.Len() == 0
		case pcommon.ValueTypeSlice:
			return c.RemoveEmptyLists && v.Slice().Len() == 0
		case pcommon.ValueTypeStr:
			return shouldFilterString(v.Str(), c.EmptyStringValues)
		}

		return false
	})
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
func cleanMetricAttrs(metric pmetric.Metric, c Config) {
	switch metric.Type() {
	case pmetric.MetricTypeGauge:
		dps := metric.Gauge().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c)
		}

	case pmetric.MetricTypeHistogram:
		dps := metric.Histogram().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c)
		}
	case pmetric.MetricTypeSum:
		dps := metric.Sum().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c)
		}
	case pmetric.MetricTypeSummary:
		dps := metric.Summary().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c)
		}
	case pmetric.MetricTypeExponentialHistogram:
		dps := metric.ExponentialHistogram().DataPoints()
		for i := 0; i < dps.Len(); i++ {
			dp := dps.At(i)
			cleanMap(dp.Attributes(), c)
		}
	default:
		// skip metric if None or unknown type
	}
}

// cleanLogBody removes empty values from the log body.
func cleanLogBody(lr plog.LogRecord, c Config) {
	body := lr.Body()
	switch body.Type() {
	case pcommon.ValueTypeMap:
		bodyMap := body.Map()
		cleanMap(bodyMap, c)
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
