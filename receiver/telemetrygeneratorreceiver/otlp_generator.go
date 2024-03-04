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

// otlpLogGenerator is a replay generator for logs, metrics, and traces.
// It generates a stream of telemetry based on the json embedded in the configuration,
// each record identical save for the timestamp.
type otlpLogGenerator struct {
	cfg           GeneratorConfig
	logger        *zap.Logger
	telemetryType component.DataType
	logs          plog.Logs
	metrics       pmetric.Metrics
	traces        ptrace.Traces
	tracesStart   time.Time
}

func newOtlpGenerator(cfg GeneratorConfig, logger *zap.Logger) *otlpLogGenerator {
	lg := &otlpLogGenerator{
		cfg:    cfg,
		logger: logger,
		logs:   plog.NewLogs(),
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
			logger.Warn("error unmarshalling otlp metrics json: %s", zap.Error(err))
		}
	case component.DataTypeTraces:
		marshaler := ptrace.JSONUnmarshaler{}
		lg.traces, err = marshaler.UnmarshalTraces(jsonBytes)
		// validation should catch this error
		if err != nil {
			logger.Warn("error unmarshalling otlp traces json: %s", zap.Error(err))
		}
		lg.tracesStart = findFirstTraceStartTime(lg.traces)
	}

	return lg
}

func findFirstTraceStartTime(traces ptrace.Traces) time.Time {
	var timeZero time.Time
	for i := 0; i < traces.ResourceSpans().Len(); i++ {
		resourceSpans := traces.ResourceSpans().At(i)
		for k := 0; k < resourceSpans.ScopeSpans().Len(); k++ {
			scopeSpans := resourceSpans.ScopeSpans().At(k)
			for j := 0; j < scopeSpans.Spans().Len(); j++ {
				span := scopeSpans.Spans().At(j)
				if span.StartTimestamp().AsTime().Before(timeZero) {
					timeZero = span.StartTimestamp().AsTime()
					continue
				}
				if timeZero.IsZero() {
					timeZero = span.StartTimestamp().AsTime()
				}
			}
		}
	}
	return timeZero
}

func (g *otlpLogGenerator) generateLogs() plog.Logs {
	for i := 0; i < g.logs.ResourceLogs().Len(); i++ {
		resourceLogs := g.logs.ResourceLogs().At(i)
		for k := 0; k < resourceLogs.ScopeLogs().Len(); k++ {
			scopeLogs := resourceLogs.ScopeLogs().At(k)
			for j := 0; j < scopeLogs.LogRecords().Len(); j++ {
				log := scopeLogs.LogRecords().At(j)
				log.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
			}
		}
	}
	return g.logs
}

func (g *otlpLogGenerator) generateMetrics() pmetric.Metrics {
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
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
					}
				case pmetric.MetricTypeGauge:
					dps := metric.Gauge().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
					}
				case pmetric.MetricTypeSummary:
					dps := metric.Summary().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
					}
				case pmetric.MetricTypeHistogram:
					dps := metric.Histogram().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
					}
				case pmetric.MetricTypeExponentialHistogram:
					dps := metric.ExponentialHistogram().DataPoints()
					for l := 0; l < dps.Len(); l++ {
						dps.At(l).SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
					}
				}
			}
		}
	}

	return g.metrics
}

// find the first trace start time, timeZero
// delta := time.Now() - timeZero
// for each span, span.StartTimestamp = span.StartTimestamp + delta
// for each span, span.EndTimestamp = span.StartTimestamp + original span length

func (g *otlpLogGenerator) generateTraces() ptrace.Traces {

	delta := time.Now().Sub(g.tracesStart)

	for i := 0; i < g.traces.ResourceSpans().Len(); i++ {
		resourceSpans := g.traces.ResourceSpans().At(i)
		for k := 0; k < resourceSpans.ScopeSpans().Len(); k++ {
			scopeSpans := resourceSpans.ScopeSpans().At(k)
			for j := 0; j < scopeSpans.Spans().Len(); j++ {
				span := scopeSpans.Spans().At(j)
				spanDuration := span.EndTimestamp().AsTime().Sub(span.StartTimestamp().AsTime())
				span.SetStartTimestamp(pcommon.NewTimestampFromTime(span.StartTimestamp().AsTime().Add(delta)))
				span.SetEndTimestamp(pcommon.NewTimestampFromTime(span.StartTimestamp().AsTime().Add(spanDuration)))
			}
		}
	}

	return g.traces
}
