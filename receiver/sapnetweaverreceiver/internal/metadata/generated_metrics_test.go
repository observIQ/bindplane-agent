// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type testMetricsSet int

const (
	testMetricsSetDefault testMetricsSet = iota
	testMetricsSetAll
	testMetricsSetNo
)

func TestMetricsBuilder(t *testing.T) {
	tests := []struct {
		name       string
		metricsSet testMetricsSet
	}{
		{
			name:       "default",
			metricsSet: testMetricsSetDefault,
		},
		{
			name:       "all_metrics",
			metricsSet: testMetricsSetAll,
		},
		{
			name:       "no_metrics",
			metricsSet: testMetricsSetNo,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start := pcommon.Timestamp(1_000_000_000)
			ts := pcommon.Timestamp(1_000_001_000)
			observedZapCore, observedLogs := observer.New(zap.WarnLevel)
			settings := receivertest.NewNopCreateSettings()
			settings.Logger = zap.New(observedZapCore)
			mb := NewMetricsBuilder(loadConfig(t, test.name), settings, WithStartTime(start))

			expectedWarnings := 0
			if test.metricsSet == testMetricsSetDefault {
				assert.Equal(t, "[WARNING] Please set `enabled` field explicitly for `sapnetweaver.work_processes.count`: This metric will be disabled by default soon.", observedLogs.All()[expectedWarnings].Message)
				expectedWarnings++
			}
			assert.Equal(t, expectedWarnings, observedLogs.Len())

			defaultMetricsCount := 0
			allMetricsCount := 0

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverAbapUpdateErrorCountDataPoint(ts, 1, AttributeControlState(1))

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverCacheEvictionsDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverCacheHitsDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverConnectionErrorCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverHostCPUUtilizationDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverHostMemoryVirtualOverheadDataPoint(ts, 1)

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverHostMemoryVirtualSwapDataPoint(ts, 1)

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverHostSpoolListUsedDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverIcmAvailabilityDataPoint(ts, 1, AttributeControlState(1))

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverJobAbortedDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverLocksEnqueueCountDataPoint(ts, 1)

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverMemoryConfiguredDataPoint(ts, 1)

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverMemoryFreeDataPoint(ts, 1)

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverMemorySwapSpaceUtilizationDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverQueueCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverQueuePeakCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverRequestCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverRequestTimeoutCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverResponseDurationDataPoint(ts, "1", AttributeResponseType(1))

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverSessionCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverSessionsBrowserCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverSessionsEjbCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverSessionsHTTPCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverSessionsSecurityCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverSessionsWebCountDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverShortDumpsRateDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverSystemAvailabilityDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverSystemUtilizationDataPoint(ts, "1")

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordSapnetweaverWorkProcessesCountDataPoint(ts, "1")

			metrics := mb.Emit(WithSapnetweaverInstance("attr-val"), WithSapnetweaverNode("attr-val"))

			if test.metricsSet == testMetricsSetNo {
				assert.Equal(t, 0, metrics.ResourceMetrics().Len())
				return
			}

			assert.Equal(t, 1, metrics.ResourceMetrics().Len())
			rm := metrics.ResourceMetrics().At(0)
			attrCount := 0
			enabledAttrCount := 0
			attrVal, ok := rm.Resource().Attributes().Get("sapnetweaver.instance")
			attrCount++
			assert.Equal(t, mb.resourceAttributesSettings.SapnetweaverInstance.Enabled, ok)
			if mb.resourceAttributesSettings.SapnetweaverInstance.Enabled {
				enabledAttrCount++
				assert.EqualValues(t, "attr-val", attrVal.Str())
			}
			attrVal, ok = rm.Resource().Attributes().Get("sapnetweaver.node")
			attrCount++
			assert.Equal(t, mb.resourceAttributesSettings.SapnetweaverNode.Enabled, ok)
			if mb.resourceAttributesSettings.SapnetweaverNode.Enabled {
				enabledAttrCount++
				assert.EqualValues(t, "attr-val", attrVal.Str())
			}
			assert.Equal(t, enabledAttrCount, rm.Resource().Attributes().Len())
			assert.Equal(t, attrCount, 2)

			assert.Equal(t, 1, rm.ScopeMetrics().Len())
			ms := rm.ScopeMetrics().At(0).Metrics()
			if test.metricsSet == testMetricsSetDefault {
				assert.Equal(t, defaultMetricsCount, ms.Len())
			}
			if test.metricsSet == testMetricsSetAll {
				assert.Equal(t, allMetricsCount, ms.Len())
			}
			validatedMetrics := make(map[string]bool)
			for i := 0; i < ms.Len(); i++ {
				switch ms.At(i).Name() {
				case "sapnetweaver.abap.update.error.count":
					assert.False(t, validatedMetrics["sapnetweaver.abap.update.error.count"], "Found a duplicate in the metrics slice: sapnetweaver.abap.update.error.count")
					validatedMetrics["sapnetweaver.abap.update.error.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The amount of ABAP errors in update.", ms.At(i).Description())
					assert.Equal(t, "", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("state")
					assert.True(t, ok)
					assert.Equal(t, "grey", attrVal.Str())
				case "sapnetweaver.cache.evictions":
					assert.False(t, validatedMetrics["sapnetweaver.cache.evictions"], "Found a duplicate in the metrics slice: sapnetweaver.cache.evictions")
					validatedMetrics["sapnetweaver.cache.evictions"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The number of evicted entries.", ms.At(i).Description())
					assert.Equal(t, "{entries}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.cache.hits":
					assert.False(t, validatedMetrics["sapnetweaver.cache.hits"], "Found a duplicate in the metrics slice: sapnetweaver.cache.hits")
					validatedMetrics["sapnetweaver.cache.hits"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "The cache hit percentage.", ms.At(i).Description())
					assert.Equal(t, "%", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.connection.error.count":
					assert.False(t, validatedMetrics["sapnetweaver.connection.error.count"], "Found a duplicate in the metrics slice: sapnetweaver.connection.error.count")
					validatedMetrics["sapnetweaver.connection.error.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The amount of connection errors.", ms.At(i).Description())
					assert.Equal(t, "", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.host.cpu.utilization":
					assert.False(t, validatedMetrics["sapnetweaver.host.cpu.utilization"], "Found a duplicate in the metrics slice: sapnetweaver.host.cpu.utilization")
					validatedMetrics["sapnetweaver.host.cpu.utilization"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "The CPU utilization percentage.", ms.At(i).Description())
					assert.Equal(t, "%", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.host.memory.virtual.overhead":
					assert.False(t, validatedMetrics["sapnetweaver.host.memory.virtual.overhead"], "Found a duplicate in the metrics slice: sapnetweaver.host.memory.virtual.overhead")
					validatedMetrics["sapnetweaver.host.memory.virtual.overhead"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Virtualization System Memory Overhead.", ms.At(i).Description())
					assert.Equal(t, "bytes", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.host.memory.virtual.swap":
					assert.False(t, validatedMetrics["sapnetweaver.host.memory.virtual.swap"], "Found a duplicate in the metrics slice: sapnetweaver.host.memory.virtual.swap")
					validatedMetrics["sapnetweaver.host.memory.virtual.swap"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Virtualization System Swap Memory.", ms.At(i).Description())
					assert.Equal(t, "bytes", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.host.spool_list.used":
					assert.False(t, validatedMetrics["sapnetweaver.host.spool_list.used"], "Found a duplicate in the metrics slice: sapnetweaver.host.spool_list.used")
					validatedMetrics["sapnetweaver.host.spool_list.used"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Host Spool List Used.", ms.At(i).Description())
					assert.Equal(t, "", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.icm_availability":
					assert.False(t, validatedMetrics["sapnetweaver.icm_availability"], "Found a duplicate in the metrics slice: sapnetweaver.icm_availability")
					validatedMetrics["sapnetweaver.icm_availability"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "ICM Availability (color value from alert tree).", ms.At(i).Description())
					assert.Equal(t, "", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("state")
					assert.True(t, ok)
					assert.Equal(t, "grey", attrVal.Str())
				case "sapnetweaver.job.aborted":
					assert.False(t, validatedMetrics["sapnetweaver.job.aborted"], "Found a duplicate in the metrics slice: sapnetweaver.job.aborted")
					validatedMetrics["sapnetweaver.job.aborted"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The amount of aborted jobs.", ms.At(i).Description())
					assert.Equal(t, "", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.locks.enqueue.count":
					assert.False(t, validatedMetrics["sapnetweaver.locks.enqueue.count"], "Found a duplicate in the metrics slice: sapnetweaver.locks.enqueue.count")
					validatedMetrics["sapnetweaver.locks.enqueue.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "Count of Enqueued Locks.", ms.At(i).Description())
					assert.Equal(t, "{locks}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.memory.configured":
					assert.False(t, validatedMetrics["sapnetweaver.memory.configured"], "Found a duplicate in the metrics slice: sapnetweaver.memory.configured")
					validatedMetrics["sapnetweaver.memory.configured"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The amount of configured memory.", ms.At(i).Description())
					assert.Equal(t, "By", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.memory.free":
					assert.False(t, validatedMetrics["sapnetweaver.memory.free"], "Found a duplicate in the metrics slice: sapnetweaver.memory.free")
					validatedMetrics["sapnetweaver.memory.free"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The amount of free memory.", ms.At(i).Description())
					assert.Equal(t, "By", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.memory.swap_space.utilization":
					assert.False(t, validatedMetrics["sapnetweaver.memory.swap_space.utilization"], "Found a duplicate in the metrics slice: sapnetweaver.memory.swap_space.utilization")
					validatedMetrics["sapnetweaver.memory.swap_space.utilization"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "The swap space utilization percentage.", ms.At(i).Description())
					assert.Equal(t, "%", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.queue.count":
					assert.False(t, validatedMetrics["sapnetweaver.queue.count"], "Found a duplicate in the metrics slice: sapnetweaver.queue.count")
					validatedMetrics["sapnetweaver.queue.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The queue length.", ms.At(i).Description())
					assert.Equal(t, "{entries}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.queue_peak.count":
					assert.False(t, validatedMetrics["sapnetweaver.queue_peak.count"], "Found a duplicate in the metrics slice: sapnetweaver.queue_peak.count")
					validatedMetrics["sapnetweaver.queue_peak.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The peak queue length.", ms.At(i).Description())
					assert.Equal(t, "{entries}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.request.count":
					assert.False(t, validatedMetrics["sapnetweaver.request.count"], "Found a duplicate in the metrics slice: sapnetweaver.request.count")
					validatedMetrics["sapnetweaver.request.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The amount of requests made.", ms.At(i).Description())
					assert.Equal(t, "", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.request.timeout.count":
					assert.False(t, validatedMetrics["sapnetweaver.request.timeout.count"], "Found a duplicate in the metrics slice: sapnetweaver.request.timeout.count")
					validatedMetrics["sapnetweaver.request.timeout.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The amount of timed out requests.", ms.At(i).Description())
					assert.Equal(t, "", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.response.duration":
					assert.False(t, validatedMetrics["sapnetweaver.response.duration"], "Found a duplicate in the metrics slice: sapnetweaver.response.duration")
					validatedMetrics["sapnetweaver.response.duration"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The response time duration.", ms.At(i).Description())
					assert.Equal(t, "ms", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("response_type")
					assert.True(t, ok)
					assert.Equal(t, "transaction", attrVal.Str())
				case "sapnetweaver.session.count":
					assert.False(t, validatedMetrics["sapnetweaver.session.count"], "Found a duplicate in the metrics slice: sapnetweaver.session.count")
					validatedMetrics["sapnetweaver.session.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The amount of of sessions created.", ms.At(i).Description())
					assert.Equal(t, "", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.sessions.browser.count":
					assert.False(t, validatedMetrics["sapnetweaver.sessions.browser.count"], "Found a duplicate in the metrics slice: sapnetweaver.sessions.browser.count")
					validatedMetrics["sapnetweaver.sessions.browser.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The number of Browser Sessions.", ms.At(i).Description())
					assert.Equal(t, "{sessions}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.sessions.ejb.count":
					assert.False(t, validatedMetrics["sapnetweaver.sessions.ejb.count"], "Found a duplicate in the metrics slice: sapnetweaver.sessions.ejb.count")
					validatedMetrics["sapnetweaver.sessions.ejb.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The number of EJB Sessions.", ms.At(i).Description())
					assert.Equal(t, "{sessions}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.sessions.http.count":
					assert.False(t, validatedMetrics["sapnetweaver.sessions.http.count"], "Found a duplicate in the metrics slice: sapnetweaver.sessions.http.count")
					validatedMetrics["sapnetweaver.sessions.http.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The number of HTTP Sessions.", ms.At(i).Description())
					assert.Equal(t, "{sessions}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.sessions.security.count":
					assert.False(t, validatedMetrics["sapnetweaver.sessions.security.count"], "Found a duplicate in the metrics slice: sapnetweaver.sessions.security.count")
					validatedMetrics["sapnetweaver.sessions.security.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The number of Security Sessions.", ms.At(i).Description())
					assert.Equal(t, "{sessions}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.sessions.web.count":
					assert.False(t, validatedMetrics["sapnetweaver.sessions.web.count"], "Found a duplicate in the metrics slice: sapnetweaver.sessions.web.count")
					validatedMetrics["sapnetweaver.sessions.web.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The number of Web Sessions.", ms.At(i).Description())
					assert.Equal(t, "{sessions}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.short_dumps.rate":
					assert.False(t, validatedMetrics["sapnetweaver.short_dumps.rate"], "Found a duplicate in the metrics slice: sapnetweaver.short_dumps.rate")
					validatedMetrics["sapnetweaver.short_dumps.rate"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The rate of Short Dumps.", ms.At(i).Description())
					assert.Equal(t, "{dumps/min}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.system.availability":
					assert.False(t, validatedMetrics["sapnetweaver.system.availability"], "Found a duplicate in the metrics slice: sapnetweaver.system.availability")
					validatedMetrics["sapnetweaver.system.availability"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "The system availability percentage.", ms.At(i).Description())
					assert.Equal(t, "%", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.system.utilization":
					assert.False(t, validatedMetrics["sapnetweaver.system.utilization"], "Found a duplicate in the metrics slice: sapnetweaver.system.utilization")
					validatedMetrics["sapnetweaver.system.utilization"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "The system utilization percentage.", ms.At(i).Description())
					assert.Equal(t, "%", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				case "sapnetweaver.work_processes.count":
					assert.False(t, validatedMetrics["sapnetweaver.work_processes.count"], "Found a duplicate in the metrics slice: sapnetweaver.work_processes.count")
					validatedMetrics["sapnetweaver.work_processes.count"] = true
					assert.Equal(t, pmetric.MetricTypeSum, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Sum().DataPoints().Len())
					assert.Equal(t, "The number of active work processes.", ms.At(i).Description())
					assert.Equal(t, "{work processes}", ms.At(i).Unit())
					assert.Equal(t, false, ms.At(i).Sum().IsMonotonic())
					assert.Equal(t, pmetric.AggregationTemporalityCumulative, ms.At(i).Sum().AggregationTemporality())
					dp := ms.At(i).Sum().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
				}
			}
		})
	}
}

func loadConfig(t *testing.T, name string) MetricsSettings {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	sub, err := cm.Sub(name)
	require.NoError(t, err)
	cfg := DefaultMetricsSettings()
	require.NoError(t, component.UnmarshalConfig(sub, &cfg))
	return cfg
}
