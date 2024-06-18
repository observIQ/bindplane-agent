package measurements

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type ThroughputMeasurementsRegistry interface {
	RegisterThroughputMeasurements(processorID string, measurements *ThroughputMeasurements)
}

// ThroughputMeasurements represents all captured throughput metrics.
// It allows for incrementing and querying the current values of throughtput metrics
type ThroughputMeasurements struct {
	logSize, metricSize, traceSize      *int64Counter
	logCount, datapointCount, spanCount *int64Counter
	attributes                          attribute.Set
}

func NewThroughputMetrics(mp metric.MeterProvider, processorID string, extraAttributes map[string]string) (*ThroughputMeasurements, error) {
	meter := mp.Meter("github.com/observiq/bindplane-agent/internal/measurements")

	logSize, err := meter.Int64Counter(
		metricName("log_data_size"),
		metric.WithDescription("Size of the log package passed to the processor"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, fmt.Errorf("create log_data_size counter: %w", err)
	}

	metricSize, err := meter.Int64Counter(
		metricName("metric_data_size"),
		metric.WithDescription("Size of the metric package passed to the processor"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, fmt.Errorf("create metric_data_size counter: %w", err)
	}

	traceSize, err := meter.Int64Counter(
		metricName("trace_data_size"),
		metric.WithDescription("Size of the trace package passed to the processor"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, fmt.Errorf("create trace_data_size counter: %w", err)
	}

	logCount, err := meter.Int64Counter(
		metricName("log_count"),
		metric.WithDescription("Count of the number log records passed to the processor"),
		metric.WithUnit("{logs}"),
	)
	if err != nil {
		return nil, fmt.Errorf("create log_count counter: %w", err)
	}

	datapointCount, err := meter.Int64Counter(
		metricName("metric_count"),
		metric.WithDescription("Count of the number datapoints passed to the processor"),
		metric.WithUnit("{datapoints}"),
	)
	if err != nil {
		return nil, fmt.Errorf("create metric_count counter: %w", err)
	}

	spanCount, err := meter.Int64Counter(
		metricName("trace_count"),
		metric.WithDescription("Count of the number spans passed to the processor"),
		metric.WithUnit("{spans}"),
	)
	if err != nil {
		return nil, fmt.Errorf("create trace_count counter: %w", err)
	}

	attrs := createMeasurementsAttributeSet(processorID, extraAttributes)

	return &ThroughputMeasurements{
		logSize:        newInt64Counter(logSize, attrs),
		logCount:       newInt64Counter(logCount, attrs),
		metricSize:     newInt64Counter(metricSize, attrs),
		datapointCount: newInt64Counter(datapointCount, attrs),
		traceSize:      newInt64Counter(traceSize, attrs),
		spanCount:      newInt64Counter(spanCount, attrs),
		attributes:     attrs,
	}, nil
}

func (tm *ThroughputMeasurements) AddLogs(ctx context.Context, l plog.Logs) {
	sizer := plog.ProtoMarshaler{}

	tm.logSize.Add(ctx, int64(sizer.LogsSize(l)))
	tm.logCount.Add(ctx, int64(l.LogRecordCount()))
}

func (tm *ThroughputMeasurements) AddMetrics(ctx context.Context, m pmetric.Metrics) {
	sizer := pmetric.ProtoMarshaler{}

	tm.metricSize.Add(ctx, int64(sizer.MetricsSize(m)))
	tm.datapointCount.Add(ctx, int64(m.DataPointCount()))
}

func (tm *ThroughputMeasurements) AddTraces(ctx context.Context, t ptrace.Traces) {
	sizer := ptrace.ProtoMarshaler{}

	tm.traceSize.Add(ctx, int64(sizer.TracesSize(t)))
	tm.spanCount.Add(ctx, int64(t.SpanCount()))
}

func (tm *ThroughputMeasurements) LogSize() int64 {
	return tm.logSize.Val()
}

func (tm *ThroughputMeasurements) MetricSize() int64 {
	return tm.metricSize.Val()
}

func (tm *ThroughputMeasurements) TraceSize() int64 {
	return tm.traceSize.Val()
}

func (tm *ThroughputMeasurements) LogCount() int64 {
	return tm.logCount.Val()
}

func (tm *ThroughputMeasurements) DatapointCount() int64 {
	return tm.datapointCount.Val()
}

func (tm *ThroughputMeasurements) SpanCount() int64 {
	return tm.spanCount.Val()
}

func (tm *ThroughputMeasurements) Attributes() attribute.Set {
	return tm.attributes
}

// int64Counter combines a metric.Int64Counter with a atomic.Int64 so that the value of the counter may be
// retrieved.
// The value of the metric counter and val are not guaranteed to be synchronized, but will be eventually consistent.
type int64Counter struct {
	counter    metric.Int64Counter
	val        atomic.Int64
	attributes attribute.Set
}

func newInt64Counter(counter metric.Int64Counter, attributes attribute.Set) *int64Counter {
	return &int64Counter{
		counter:    counter,
		attributes: attributes,
	}
}

func (i *int64Counter) Add(ctx context.Context, delta int64) {
	i.counter.Add(ctx, delta, metric.WithAttributeSet(i.attributes))
	i.val.Add(delta)
}

func (i *int64Counter) Val() int64 {
	return i.val.Load()
}

func metricName(metric string) string {
	return fmt.Sprintf("otelcol_processor_throughputmeasurement_%s", metric)
}

func createMeasurementsAttributeSet(processorID string, extraAttributes map[string]string) attribute.Set {
	attrs := make([]attribute.KeyValue, 0, len(extraAttributes)+1)

	attrs = append(attrs, attribute.String("processor", processorID))
	for k, v := range extraAttributes {
		attrs = append(attrs, attribute.String(k, v))
	}

	return attribute.NewSet(attrs...)
}

type ConcreteThroughputMeasurementsRegistry struct {
	measurements     *sync.Map
	emitCountMetrics bool
}

func NewConcreteThroughputMeasurementsRegistry(emitCountMetrics bool) *ConcreteThroughputMeasurementsRegistry {
	return &ConcreteThroughputMeasurementsRegistry{
		measurements:     &sync.Map{},
		emitCountMetrics: emitCountMetrics,
	}
}

func (ctmr *ConcreteThroughputMeasurementsRegistry) RegisterThroughputMeasurements(processorID string, measurements *ThroughputMeasurements) {
	ctmr.measurements.Store(processorID, measurements)
}

func (ctmr *ConcreteThroughputMeasurementsRegistry) Reset() {
	ctmr.measurements = &sync.Map{}
}

func (ctmr *ConcreteThroughputMeasurementsRegistry) OTLPMeasurements() pmetric.Metrics {
	m := pmetric.NewMetrics()
	rm := m.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()

	ctmr.measurements.Range(func(processor, value any) bool {
		tm := value.(*ThroughputMeasurements)
		OTLPThroughputMeasurements(tm, ctmr.emitCountMetrics).MoveAndAppendTo(sm.Metrics())
		return true
	})

	return m
}
