package snowflakeexporter

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/database"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/metrics"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/utility"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

type metricModel interface {
	AddMetric(r pmetric.ResourceMetrics, s pmetric.ScopeMetrics, m pmetric.Metric)
	BatchInsert(ctx context.Context, db database.Database) error
}

type metricsExporter struct {
	cfg    *Config
	logger *zap.Logger
	db     database.Database
	models map[string]metricModel
}

func newMetricsExporter(
	ctx context.Context,
	c *Config,
	params exporter.CreateSettings,
	newDatabase func(ctx context.Context, dsn string) (database.Database, error),
) (*metricsExporter, error) {
	dsn := utility.CreateDSN(
		c.Username,
		c.Password,
		c.AccountIdentifier,
		c.Database,
	)

	db, err := newDatabase(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create new database connection for metrics: %w", err)
	}

	return &metricsExporter{
		cfg:    c,
		logger: params.Logger,
		db:     db,
		models: map[string]metricModel{},
	}, nil
}

func (me *metricsExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (me *metricsExporter) start(ctx context.Context, _ component.Host) error {
	err := me.db.CreateSchema(ctx, me.cfg.Metrics.Schema)
	if err != nil {
		return fmt.Errorf("failed to create metrics schema: %w", err)
	}

	// init metric models
	me.models["sums"] = metrics.NewSumModel(me.logger, me.cfg.Warehouse, me.cfg.Metrics.Schema, me.cfg.Metrics.Table)
	me.models["gauges"] = metrics.NewGaugeModel(me.logger, me.cfg.Warehouse, me.cfg.Metrics.Schema, me.cfg.Metrics.Table)
	me.models["histograms"] = metrics.NewHistogramModel(me.logger, me.cfg.Warehouse, me.cfg.Metrics.Schema, me.cfg.Metrics.Table)
	me.models["exponentialHistograms"] = metrics.NewExponentialHistogramModel(me.logger, me.cfg.Warehouse, me.cfg.Metrics.Schema, me.cfg.Metrics.Table)
	me.models["summaries"] = metrics.NewSummaryModel(me.logger, me.cfg.Warehouse, me.cfg.Metrics.Schema, me.cfg.Metrics.Table)

	// create metric tables
	err = me.db.CreateTable(ctx, me.cfg.Database, me.cfg.Metrics.Schema, me.cfg.Metrics.Table, metrics.CreateSumMetricTableTemplate)
	if err != nil {
		return fmt.Errorf("failed to create sum metrics table: %w", err)
	}
	err = me.db.CreateTable(ctx, me.cfg.Database, me.cfg.Metrics.Schema, me.cfg.Metrics.Table, metrics.CreateGaugeMetricTableTemplate)
	if err != nil {
		return fmt.Errorf("failed to create gauge metrics table: %w", err)
	}
	err = me.db.CreateTable(ctx, me.cfg.Database, me.cfg.Metrics.Schema, me.cfg.Metrics.Table, metrics.CreateHistogramMetricTableTemplate)
	if err != nil {
		return fmt.Errorf("failed to create histogram metrics table: %w", err)
	}
	err = me.db.CreateTable(ctx, me.cfg.Database, me.cfg.Metrics.Schema, me.cfg.Metrics.Table, metrics.CreateExponentialHistogramMetricTableTemplate)
	if err != nil {
		return fmt.Errorf("failed to create exponential histogram metrics table: %w", err)
	}
	err = me.db.CreateTable(ctx, me.cfg.Database, me.cfg.Metrics.Schema, me.cfg.Metrics.Table, metrics.CreateSummaryMetricTableTemplate)
	if err != nil {
		return fmt.Errorf("failed to create summary metrics table: %w", err)
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

	// loop through metrics and add to corresponding metric model
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		resourceMetric := md.ResourceMetrics().At(i)

		for j := 0; j < resourceMetric.ScopeMetrics().Len(); j++ {
			scopeMetric := resourceMetric.ScopeMetrics().At(j)

			for k := 0; k < scopeMetric.Metrics().Len(); k++ {
				metric := scopeMetric.Metrics().At(k)

				switch metric.Type() {
				case pmetric.MetricTypeSum:
					me.models["sums"].AddMetric(resourceMetric, scopeMetric, metric)
				case pmetric.MetricTypeGauge:
					me.models["gauges"].AddMetric(resourceMetric, scopeMetric, metric)
				case pmetric.MetricTypeHistogram:
					me.models["histograms"].AddMetric(resourceMetric, scopeMetric, metric)
				case pmetric.MetricTypeExponentialHistogram:
					me.models["exponentialHistograms"].AddMetric(resourceMetric, scopeMetric, metric)
				case pmetric.MetricTypeSummary:
					me.models["summaries"].AddMetric(resourceMetric, scopeMetric, metric)
				default:
					me.logger.Warn("unsupported metric type", zap.String("type", metric.Type().String()))
				}
			}
		}
	}

	// call BatchInsert for all metric models
	wg := &sync.WaitGroup{}
	errorChan := make(chan error, len(me.models))
	for _, v := range me.models {
		wg.Add(1)
		go func(m metricModel) {
			defer wg.Done()
			errorChan <- m.BatchInsert(ctx, me.db)
		}(v)
	}
	wg.Wait()
	close(errorChan)

	// return any errors
	var errs error
	for e := range errorChan {
		errors.Join(errs, e)
	}

	return errs
}
