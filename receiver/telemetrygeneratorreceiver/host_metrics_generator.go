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
	"math"
	"math/rand"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

// hostMetricsGenerator is a generator for host metrics. It generates a sampling of host metrics
// emulating the Host Metrics receiver: github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver
type hostMetricsGenerator struct {
	otlpGenerator
	// dataPointUpdaters is a list of functions that update the data points with new values
	// for each batch of metrics that are generated
	dataPointUpdaters []func()
}

func newHostMetricsGenerator(cfg GeneratorConfig, logger *zap.Logger) metricGenerator {

	g := &hostMetricsGenerator{
		otlpGenerator: otlpGenerator{
			cfg:     cfg,
			logger:  logger,
			logs:    plog.NewLogs(),
			metrics: pmetric.NewMetrics(),
			traces:  ptrace.NewTraces(),
		},
	}
	// Load up Resources attributes

	newResource := g.metrics.ResourceMetrics().AppendEmpty()
	newResource.Resource().Attributes().FromRaw(cfg.ResourceAttributes)

	metrics := cfg.AdditionalConfig["metrics"].([]any)
	newScope := newResource.ScopeMetrics().AppendEmpty()
	for _, m := range metrics {
		metric := m.(map[string]any)
		var attributes map[string]any
		// attributes are optional
		if attr, ok := metric["attributes"]; ok {
			attributes = attr.(map[string]any)
		}

		metricType := metric["type"].(string)

		name := metric["name"].(string)
		valueMin := metric["value_min"].(int)
		valueMax := metric["value_max"].(int)
		unit := metric["unit"].(string)

		newMetric := newScope.Metrics().AppendEmpty()
		newMetric.SetUnit(unit)
		newMetric.SetName(name)

		// all the host metrics are either sums or gauges
		switch metricType {
		case "Gauge":
			newMetric.SetEmptyGauge()
			dp := newMetric.Gauge().DataPoints().AppendEmpty()
			dp.SetTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
			dp.SetStartTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
			dp.Attributes().FromRaw(attributes)

			// All the host metric Gauges are float64
			g.dataPointUpdaters = append(g.dataPointUpdaters, func() { dp.SetDoubleValue(getRandomFloat64(valueMin, valueMax)) })
		case "Sum":
			newMetric.SetEmptySum()
			dp := newMetric.Sum().DataPoints().AppendEmpty()
			dp.SetTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
			dp.SetStartTimestamp(pcommon.NewTimestampFromTime(getCurrentTime()))
			dp.Attributes().FromRaw(attributes)
			switch unit {
			case "s":
				g.dataPointUpdaters = append(g.dataPointUpdaters, func() { dp.SetDoubleValue(float64(math.Trunc(getRandomFloat64(valueMin, valueMax)*100)) / 100) })
			default:
				g.dataPointUpdaters = append(g.dataPointUpdaters, func() { dp.SetIntValue(int64(getRandomFloat64(valueMin, valueMax))) })
			}
		}
	}

	g.adjustMetricTimes()

	return g
}

// this is a variable so that it can be overridden in tests
var getRandomFloat64 = func(value_min, value_max int) float64 {
	return rand.Float64()*(float64(value_max)-float64(value_min)) + float64(value_min)
}

func (g *hostMetricsGenerator) generateMetrics() pmetric.Metrics {
	// Update the data points with new values
	for _, updater := range g.dataPointUpdaters {
		updater()
	}
	return g.otlpGenerator.generateMetrics()
}
