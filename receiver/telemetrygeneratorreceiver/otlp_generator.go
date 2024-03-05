// Copyright observIQ, Inc.
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

package telemetrygeneratorreceiver //import "github.com/observiq/bindplane-agent/receiver/telemetrygeneratorreceiver"

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// otlpGenerator is a replay generator for logs, metrics, and traces.
// It generates a stream of telemetry based on the json embedded in the configuration,
// each record identical save for the timestamp.
type otlpGenerator struct {
	cfg           GeneratorConfig
	logger        *zap.Logger
	telemetryType component.DataType
	logs          plog.Logs
	metrics       pmetric.Metrics
	traces        ptrace.Traces
	tracesStart   time.Time
}

func newOtlpGenerator(cfg GeneratorConfig, logger *zap.Logger) *otlpGenerator {
	lg := &otlpGenerator{
		cfg:     cfg,
		logger:  logger,
		logs:    plog.NewLogs(),
		metrics: pmetric.NewMetrics(),
		traces:  ptrace.NewTraces(),
	}

	// validation already proves this exists, is a string, and a component.DataType
	lg.telemetryType = component.Type(lg.cfg.AdditionalConfig["telemetry_type"].(string))
	jsonBytes := []byte(lg.cfg.AdditionalConfig["otlp_json"].(string))

	var err error
	switch lg.telemetryType {
	case component.DataTypeLogs:
		marshaler := plog.JSONUnmarshaler{}
		lg.logs, err = marshaler.UnmarshalLogs(jsonBytes)
		// validation should catch this error
		if err != nil {
			logger.Warn("error unmarshalling otlp logs json", zap.Error(err))
		}
	case component.DataTypeMetrics:
		marshaler := pmetric.JSONUnmarshaler{}
		lg.metrics, err = marshaler.UnmarshalMetrics(jsonBytes)
		// validation should catch this error
		if err != nil {
			logger.Warn("error unmarshalling otlp metrics json", zap.Error(err))
		}
	case component.DataTypeTraces:
		marshaler := ptrace.JSONUnmarshaler{}
		lg.traces, err = marshaler.UnmarshalTraces(jsonBytes)
		// validation should catch this error
		if err != nil {
			logger.Warn("error unmarshalling otlp traces json", zap.Error(err))
		}
		lastEndTime := findLastTraceEndTime(lg.traces)
		lg.adjustTraceTimes(lastEndTime)
	}

	return lg
}

// getCurrentTime is a variable that holds the current time function. It is used to mock time in tests.
var getCurrentTime = func() time.Time { return time.Now().UTC() }

// findLastTraceEndTime finds the span with the last end time
func findLastTraceEndTime(traces ptrace.Traces) time.Time {
	var maxTime time.Time
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		resourceSpans := traces.ResourceSpans().At(i)
		for k := 0; k < resourceSpans.ScopeSpans().Len(); k++ {
			scopeSpans := resourceSpans.ScopeSpans().At(k)
			for j := 0; j < scopeSpans.Spans().Len(); j++ {
				span := scopeSpans.Spans().At(j)
				if span.EndTimestamp().AsTime().After(maxTime) {
					maxTime = span.EndTimestamp().AsTime()
					continue
				}
				if maxTime.IsZero() {
					maxTime = span.EndTimestamp().AsTime()
				}
			}
		}
	}
	return maxTime
}

// adjustTraceTimes changes the start and end times of all spans to be relative to the current time, placing
// the span that ends at maxTime at the current time.
func (g *otlpGenerator) adjustTraceTimes(maxTime time.Time) {
	now := getCurrentTime()
	for i := 0; i < g.traces.ResourceSpans().Len(); i++ {
		resourceSpans := g.traces.ResourceSpans().At(i)
		for k := 0; k < resourceSpans.ScopeSpans().Len(); k++ {
			scopeSpans := resourceSpans.ScopeSpans().At(k)
			for j := 0; j < scopeSpans.Spans().Len(); j++ {
				span := scopeSpans.Spans().At(j)
				// delta is the differnce between this spans end time and the max end time
				delta := maxTime.Sub(span.EndTimestamp().AsTime())
				// spanDuration is the length of the span
				spanDuration := span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime())
				endTime := now
				// move each span's end time by delta
				span.SetEndTimestamp(pcommon.NewTimestampFromTime(endTime.Add(delta)))
				// set the start time to be the end time minus the original span duration
				span.SetStartTimestamp(pcommon.NewTimestampFromTime(endTime.Add(-spanDuration)))
			}
		}
	}
	// save the current time we used as a baseline to adjust the spans
	g.tracesStart = now
}

func (g *otlpGenerator) generateLogs() plog.Logs {
	for i := 0; i < g.logs.ResourceLogs().Len(); i++ {
		resourceLogs := g.logs.ResourceLogs().At(i)
		for k := 0; k < resourceLogs.ScopeLogs().Len(); k++ {
			scopeLogs := resourceLogs.ScopeLogs().At(k)
			for j := 0; j < scopeLogs.LogRecords().Len(); j++ {
				log := scopeLogs.LogRecords().At(j)
				log.SetTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
			}
		}
	}
	return g.logs
}

func (g *otlpGenerator) generateMetrics() pmetric.Metrics {
	for i := 0; i < g.metrics.ResourceMetrics().Len(); i++ {
		resourceMetrics := g.metrics.ResourceMetrics().At(i)
		for k := 0; k < resourceMetrics.ScopeMetrics().Len(); k++ {
			scopeMetrics := resourceMetrics.ScopeMetrics().At(k)
			for j := 0; j < scopeMetrics.Metrics().Len(); j++ {
				metric := scopeMetrics.Metrics().At(j)
				switch metric.Type() {
				case pmetric.MetricTypeSum:
					dps := metric.Sum().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
					}
				case pmetric.MetricTypeGauge:
					dps := metric.Gauge().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
					}
				case pmetric.MetricTypeSummary:
					dps := metric.Summary().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
					}
				case pmetric.MetricTypeHistogram:
					dps := metric.Histogram().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
					}
				case pmetric.MetricTypeExponentialHistogram:
					dps := metric.ExponentialHistogram().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
					}
				}
			}
		}
	}

	return g.metrics
}

func (g *otlpGenerator) generateTraces() ptrace.Traces {
	// calculate the time since the last baseline time we used to adjust the spans
	now := getCurrentTime()
	timeSince := now.Sub(g.tracesStart)
	// add the time since to the start and end times of all spans
	for i := 0; i < g.traces.ResourceSpans().Len(); i++ {
		resourceSpans := g.traces.ResourceSpans().At(i)
		for k := 0; k < resourceSpans.ScopeSpans().Len(); k++ {
			scopeSpans := resourceSpans.ScopeSpans().At(k)
			for j := 0; j < scopeSpans.Spans().Len(); j++ {
				span := scopeSpans.Spans().At(j)

				startTime := span.StartTimestamp().AsTime().Add(timeSince)
				span.SetStartTimestamp(pcommon.NewTimestampFromTime(startTime))

				endTime := span.EndTimestamp().AsTime().Add(timeSince)
				span.SetEndTimestamp(pcommon.NewTimestampFromTime(endTime))
			}
		}
	}
	// update the baseline time to the current time
	g.tracesStart = now

	return g.traces
}
