package snowflakeexporter

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/metrics"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/utility"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type metricModel interface {
	AddMetric(r pmetric.ResourceMetrics, s pmetric.ScopeMetrics, m pmetric.Metric, d any)
	BatchInsert(ctx context.Context) error
}

type metricsExporter struct {
	cfg    *Config
	logger *zap.Logger
	db     *sqlx.DB
	models map[string]metricModel
}

func newMetricsExporter(c *Config, params exporter.CreateSettings) (*metricsExporter, error) {
	return &metricsExporter{
		cfg:    c,
		logger: params.Logger,
		models: map[string]metricModel{},
	}, nil
}

func (me *metricsExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (me *metricsExporter) start(ctx context.Context, _ component.Host) error {
	dsn := utility.CreateDSN(
		me.cfg.Username,
		me.cfg.Password,
		me.cfg.AccountIdentifier,
		me.cfg.Database,
		me.cfg.Metrics.Schema,
	)
	db, err := utility.CreateDB(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to create new db connection for metrics: %w", err)
	}
	me.db = db

	// TODO: init more metric models
	me.models["sums"] = metrics.NewSumModel(me.logger, me.db, me.cfg.Warehouse, me.cfg.Metrics.Schema, me.cfg.Metrics.Table)
	me.models["gauges"] = metrics.NewGaugeModel(me.logger, me.db, me.cfg.Warehouse, me.cfg.Metrics.Schema, me.cfg.Metrics.Table)

	// TODO: create more metric tables
	err = utility.CreateTable(ctx, me.db, me.cfg.Database, me.cfg.Metrics.Schema, me.cfg.Metrics.Table, metrics.CreateSumMetricTableTemplate)
	if err != nil {
		return fmt.Errorf("failed to create sum metrics table: %w", err)
	}

	err = utility.CreateTable(ctx, me.db, me.cfg.Database, me.cfg.Metrics.Schema, me.cfg.Metrics.Table, metrics.CreateGaugeMetricTableTemplate)
	if err != nil {
		return fmt.Errorf("failed to create gauge metrics table: %w", err)
	}
	return nil
}

func (me *metricsExporter) shutdown(_ context.Context) error {
	if me.db != nil {
		return me.db.Close()
	}
	return nil
}

func (me *metricsExporter) metricsDataPusher(ctx context.Context, md pmetric.Metrics) error {
	me.logger.Debug("begin metricsDataPusher")

	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		resourceMetric := md.ResourceMetrics().At(i)

		for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
			scopeMetric := resourceMetric.ScopeMetrics().At(j)

			for k := 0; k < scopeMetric.Metrics().Len(); k++ {
				metric := scopeMetric.Metrics().At(k)

				// TODO: add more metrics types
				switch metric.Type() {
				case pmetric.MetricTypeSum:
					me.models["sums"].AddMetric(resourceMetric, scopeMetric, metric, metric.Sum())
				case pmetric.MetricTypeGauge:
					me.models["gauges"].AddMetric(resourceMetric, scopeMetric, metric, metric.Gauge())
				default:
					me.logger.Warn("unsupported metric type", zap.String("type", metric.Type().String()))
				}
			}
		}
	}

	wg := &sync.WaitGroup{}
	errorChan := make(chan error, len(me.models))
	for _, v := range me.models {
		wg.Add(1)
		go func(m metricModel) {
			defer wg.Done()
			errorChan <- m.BatchInsert(ctx)
		}(v)
	}
	wg.Wait()
	close(errorChan)

	var errs error
	for e := range errorChan {
		errors.Join(errs, e)
	}

	return errs
}
