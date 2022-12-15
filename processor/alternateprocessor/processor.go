package alternateprocessor

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/observiq/observiq-otel-collector/receiver/routereceiver"
)

type alternateProcessor struct {
	cfg    *Config
	logger *zap.Logger

	metricsTracker  *RollingAverage
	metricsRate     *Rate
	metricsConsumer *consumer.Metrics
	metricsSizer    pmetric.Sizer

	logsConsumer *consumer.Logs
	logsRate     *Rate
	logsTracker  *RollingAverage
	logsSizer    plog.Sizer

	tracesConsumer *consumer.Traces
	tracesRate     *Rate
	tracesTracker  *RollingAverage
	tracesSizer    ptrace.Sizer
}

type alternateProcessorOption interface {
	alternateProcessorOptionFunc(*alternateProcessor)
}

var _ alternateProcessorOption = (*alternateProcessorOptionFunc)(nil)

type alternateProcessorOptionFunc func(*alternateProcessor)

func (apo alternateProcessorOptionFunc) alternateProcessorOptionFunc(ap *alternateProcessor) {
	apo(ap)
}

func withLogsConsumer(c consumer.Logs) alternateProcessorOption {
	return alternateProcessorOptionFunc(func(ap *alternateProcessor) {
		ap.logger.Info("anotha one")
		ap.logsConsumer = &c
	})
}

func withMetricsConsumer(c consumer.Metrics) alternateProcessorOption {
	return alternateProcessorOptionFunc(func(ap *alternateProcessor) {
		ap.metricsConsumer = &c
	})
}

func withTracesConsumer(c consumer.Traces) alternateProcessorOption {
	return alternateProcessorOptionFunc(func(ap *alternateProcessor) {
		ap.tracesConsumer = &c
	})
}

func newProcessor(
	cfg *Config,
	logger *zap.Logger,
	options ...alternateProcessorOption,
) (*alternateProcessor, error) {
	ap := &alternateProcessor{
		cfg:          cfg,
		logger:       logger,
		logsSizer:    &plog.ProtoMarshaler{},
		metricsSizer: &pmetric.ProtoMarshaler{},
		tracesSizer:  &ptrace.ProtoMarshaler{},
	}

	for _, o := range options {
		o.alternateProcessorOptionFunc(ap)
	}

	var errs error
	errs = multierr.Append(errs, ap.setupLogsTracker())
	errs = multierr.Append(errs, ap.setupMetricsTracker())
	errs = multierr.Append(errs, ap.setuptracesTracker())

	return ap, errs
}

func (ap *alternateProcessor) Start(_ context.Context, _ component.Host) error {
	if ap.logsTracker != nil {
		ap.logsTracker.Start()
	}
	if ap.metricsTracker != nil {
		ap.metricsTracker.Start()
	}
	return nil
}

func (ap *alternateProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (ap *alternateProcessor) Shutdown(_ context.Context) error {
	if ap.logsTracker != nil {
		ap.logsTracker.Stop()
	}
	if ap.metricsTracker != nil {
		ap.metricsTracker.Stop()
	}
	return nil
}

func (ap *alternateProcessor) ConsumeLogs(ctx context.Context, pl plog.Logs) error {
	if !ap.cfg.Logs.Enabled {
		return ap.normalConsumeLogs(ctx, pl)
	}

	byteSize := float64(ap.logsSizer.LogsSize(pl))
	ap.logsTracker.AddBytes(byteSize)

	currentRate := ap.logsTracker.NormalizedRateValue()

	// normal case route to the original pipeline
	if currentRate <= ap.logsRate.Value {
		return ap.normalConsumeLogs(ctx, pl)
	}

	// otherwise route to the alternate pipeline
	ap.logger.Info(
		"exceeded data limit for logs, sending logs to alternate route",
		zap.Float64("currentRate", currentRate),
		zap.Float64("configuredRate", ap.logsRate.Value),
		zap.String("route", ap.cfg.Logs.Route),
	)

	err := routereceiver.RouteLogs(ctx, ap.cfg.Logs.Route, pl)
	if err != nil {
		ap.logger.Error("failed to send logs to alternate route", zap.Error(err))
	}
	return err
}

func (ap *alternateProcessor) ConsumeMetrics(ctx context.Context, pm pmetric.Metrics) error {
	if !ap.cfg.Metrics.Enabled {
		return ap.normalConsumeMetrics(ctx, pm)
	}

	byteSize := float64(ap.metricsSizer.MetricsSize(pm))
	ap.metricsTracker.AddBytes(byteSize)
	currentRate := ap.metricsTracker.NormalizedRateValue()

	// normal case route to the original pipeline
	if currentRate <= ap.metricsRate.Value {
		return ap.normalConsumeMetrics(ctx, pm)
	}

	// otherwise route to the alternate pipeline
	ap.logger.Info(
		"exceeded data limit for metrics, sending metrics to alternate route",
		zap.Float64("currentRate", currentRate),
		zap.Float64("configuredRate", ap.metricsRate.Value),
		zap.String("route", ap.cfg.Metrics.Route),
	)

	err := routereceiver.RouteMetrics(ctx, ap.cfg.Metrics.Route, pm)
	if err != nil {
		ap.logger.Error("failed to send metrics to alternate route", zap.Error(err))
	}
	return err
}

func (ap *alternateProcessor) ConsumeTraces(ctx context.Context, pt ptrace.Traces) error {
	if !ap.cfg.Traces.Enabled {
		ap.normalConsumeTraces(ctx, pt)
	}

	byteSize := float64(ap.tracesSizer.TracesSize(pt))
	ap.tracesTracker.AddBytes(byteSize)
	currentRate := ap.tracesTracker.NormalizedRateValue()

	if currentRate <= ap.tracesRate.Value {
		ap.normalConsumeTraces(ctx, pt)
	}

	// otherwise route to the alternate pipeline
	ap.logger.Info(
		"exceeded data limit for traces, sending traces to alternate route",
		zap.Float64("currentRate", currentRate),
		zap.Float64("configuredRate", ap.tracesRate.Value),
		zap.String("route", ap.cfg.Traces.Route),
	)

	err := routereceiver.RouteTraces(ctx, ap.cfg.Traces.Route, pt)
	if err != nil {
		ap.logger.Error("failed to send traces to alternate route", zap.Error(err))
	}
	return err
}

func (ap *alternateProcessor) setupLogsTracker() error {
	if !ap.cfg.Logs.Enabled {
		return nil
	}
	rate, err := ParseRate(ap.cfg.Logs.Rate)
	if err != nil {
		return err
	}
	ap.logsRate = rate

	numBuckets := int(ap.cfg.Logs.AggregationInterval / rate.Time.Value)
	lt, err := NewRollingAverage(numBuckets, ap.cfg.Logs.AggregationInterval)
	if err != nil {
		return err
	}

	ap.logsTracker = lt
	return nil
}

func (ap *alternateProcessor) setupMetricsTracker() error {
	if !ap.cfg.Metrics.Enabled {
		return nil
	}

	rate, err := ParseRate(ap.cfg.Metrics.Rate)
	if err != nil {
		return err
	}
	ap.metricsRate = rate

	numBuckets := int(ap.cfg.Metrics.AggregationInterval / rate.Time.Value)
	mt, err := NewRollingAverage(numBuckets, ap.cfg.Metrics.AggregationInterval)
	if err != nil {
		return err
	}

	ap.metricsTracker = mt
	return nil
}

func (ap *alternateProcessor) setuptracesTracker() error {
	if !ap.cfg.Traces.Enabled {
		return nil
	}
	rate, err := ParseRate(ap.cfg.Traces.Rate)
	if err != nil {
		return err
	}
	ap.tracesRate = rate

	numBuckets := int(ap.cfg.Traces.AggregationInterval / rate.Time.Value)
	tt, err := NewRollingAverage(numBuckets, ap.cfg.Traces.AggregationInterval)
	if err != nil {
		return err
	}
	ap.tracesTracker = tt
	return nil
}

func (ap *alternateProcessor) normalConsumeLogs(ctx context.Context, pl plog.Logs) error {
	if ap.logsConsumer == nil {
		return component.ErrNilNextConsumer
	}
	lc := *ap.logsConsumer
	return lc.ConsumeLogs(ctx, pl)
}

func (ap *alternateProcessor) normalConsumeMetrics(ctx context.Context, pm pmetric.Metrics) error {
	if ap.metricsConsumer == nil {
		return component.ErrNilNextConsumer
	}
	mc := *ap.metricsConsumer
	return mc.ConsumeMetrics(ctx, pm)
}

func (ap *alternateProcessor) normalConsumeTraces(ctx context.Context, pt ptrace.Traces) error {
	if ap.tracesConsumer == nil {
		return component.ErrNilNextConsumer
	}
	tc := *ap.tracesConsumer
	return tc.ConsumeTraces(ctx, pt)
}
