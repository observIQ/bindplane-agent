package alternateprocessor

import (
	"context"
	"sync"

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
	mux    *sync.Mutex

	metricsTracker  *RateTracker
	metricsConsumer *consumer.Metrics
	metricsSizer    pmetric.Sizer

	logsConsumer *consumer.Logs
	logsTracker  *RateTracker
	logsSizer    plog.Sizer

	tracesSizer ptrace.Sizer
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

func newProcessor(
	cfg *Config,
	logger *zap.Logger,
	options ...alternateProcessorOption) (*alternateProcessor, error) {
	ap := &alternateProcessor{
		cfg:          cfg,
		logger:       logger,
		logsSizer:    &plog.ProtoMarshaler{},
		metricsSizer: &pmetric.ProtoMarshaler{},
		tracesSizer:  &ptrace.ProtoMarshaler{},
		mux:          &sync.Mutex{},
	}

	for _, o := range options {
		o.alternateProcessorOptionFunc(ap)
	}

	var errs error
	errs = multierr.Append(errs, ap.setupLogsTracker())
	errs = multierr.Append(errs, ap.setupMetricsTracker())

	return ap, nil
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
	byteSize := int64(ap.logsSizer.LogsSize(pl))
	ap.logsTracker.Add(byteSize)

	currentRate := ap.logsTracker.GetRate()

	// normal case route to the original pipeline
	if currentRate <= ap.cfg.Logs.Limit {
		if ap.logsConsumer == nil {
			return component.ErrNilNextConsumer
		}
		lc := *ap.logsConsumer
		return lc.ConsumeLogs(ctx, pl)
	}

	// otherwise route to the alternate pipeline
	ap.logger.Info(
		"exceeded limit for logs, sending logs to alternate route",
		zap.Float64("currentRate", currentRate),
		zap.Float64("configuredRate", ap.cfg.Logs.Limit),
		zap.String("route", ap.cfg.Logs.Route),
	)

	err := routereceiver.RouteLogs(ctx, ap.cfg.Logs.Route, pl)
	if err != nil {
		ap.logger.Error("failed to send logs to alternate route", zap.Error(err))
	}
	return err
}

func (ap *alternateProcessor) ConsumeMetrics(ctx context.Context, pm pmetric.Metrics) error {
	byteSize := int64(ap.metricsSizer.MetricsSize(pm))
	ap.metricsTracker.Add(byteSize)

	currentRate := ap.logsTracker.GetRate()

	// normal case route to the original pipeline
	if currentRate <= ap.cfg.Logs.Limit {
		if ap.metricsConsumer == nil {
			return component.ErrNilNextConsumer
		}
		mc := *ap.metricsConsumer
		return mc.ConsumeMetrics(ctx, pm)
	}

	// otherwise route to the alternate pipeline
	ap.logger.Info(
		"exceeded limit for logs, sending logs to alternate route",
		zap.Float64("currentRate", currentRate),
		zap.Float64("configuredRate", ap.cfg.Logs.Limit),
		zap.String("route", ap.cfg.Logs.Route),
	)

	err := routereceiver.RouteMetrics(ctx, ap.cfg.Metrics.Route, pm)
	if err != nil {
		ap.logger.Error("failed to send logs to alternate route", zap.Error(err))
	}
	return err
}

func (ap *alternateProcessor) setupLogsTracker() error {
	if ap.cfg.Logs == nil {
		return nil
	}
	lt, err := ap.cfg.Logs.Rate.Build()
	if err != nil {
		return err
	}
	ap.logsTracker = lt
	return nil
}

func (ap *alternateProcessor) setupMetricsTracker() error {
	if ap.cfg.Metrics == nil {
		return nil
	}
	mt, err := ap.cfg.Metrics.Rate.Build()
	if err != nil {
		return err
	}
	ap.metricsTracker = mt
	return nil
}
