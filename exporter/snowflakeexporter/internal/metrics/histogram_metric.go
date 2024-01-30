package metrics

import (
	"context"
	"fmt"

	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/database"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/utility"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

const (
	CreateHistogramMetricTableTemplate = `
	CREATE TABLE IF NOT EXISTS "%s"."%s_histogram" (
		"ResourceSchemaURL" VARCHAR,
		"ResourceDroppedAttributesCount" INT,
		"ResourceAttributes" VARCHAR,
		"ScopeSchemaURL" VARCHAR,
		"ScopeName" VARCHAR,
		"ScopeVersion" VARCHAR,
		"ScopeDroppedAttributesCount" INT,
		"ScopeAttributes" VARCHAR,
		"MetricName" VARCHAR,
		"MetricDescription" VARCHAR,
		"MetricUnit" VARCHAR,
		"AggregationTemporality" VARCHAR,
		"Attributes" VARCHAR,
		"StartTimestamp" TIMESTAMP_NTZ,
		"Timestamp" TIMESTAMP_NTZ,
		"Count" INT,
		"Sum" NUMBER,
		"Flags" INT,
		"Min" INT,
		"Max" INT,
		"BucketCounts" VARCHAR,
		"ExplicitBounds" VARCHAR,
		"ExemplarAttributes" VARCHAR,
		"ExemplarTimestamps" VARCHAR,
		"ExemplarTraceIDs" VARCHAR,
		"ExemplarSpanIDs" VARCHAR,
		"ExemplarValues" VARCHAR
	);`

	insertIntoHistogramMeticTableTemplate = `
	INSERT INTO "%s"."%s_histogram" (
		"ResourceSchemaURL",
		"ResourceDroppedAttributesCount",
		"ResourceAttributes",
		"ScopeSchemaURL",
		"ScopeName",
		"ScopeVersion",
		"ScopeDroppedAttributesCount",
		"ScopeAttributes",
		"MetricName",
		"MetricDescription",
		"MetricUnit",
		"Attributes",
		"StartTimestamp",
		"Timestamp",
		"Count",
		"Sum",
		"Flags",
		"Min",
		"Max",
		"BucketCounts",
		"ExplicitBounds",
		"ExemplarAttributes",
		"ExemplarTimestamps",
		"ExemplarTraceIDs",
		"ExemplarSpanIDs",
		"ExemplarValues"
	) VALUES (
		:rSchema,
		:rDroppedCount,
		:rAttributes,
		:sSchema,
		:sName,
		:sVersion,
		:sDroppedCount,
		:sAttributes,
		:mName,
		:mDescription,
		:mUnit,
		:attributes,
		:startTimestamp,
		:timestamp,
		:count,
		:sum,
		:flags,
		:min,
		:max,
		:bucketCounts,
		:explicitBounds,
		:eAttributes,
		:eTimestamps,
		:eTraceIDs,
		:eSpanIDs,
		:eValues
	);`
)

type HistogramModel struct {
	logger     *zap.Logger
	histograms []*histogramData
	warehouse  string
	insertSQL  string
}

type histogramData struct {
	resource  pmetric.ResourceMetrics
	scope     pmetric.ScopeMetrics
	metric    pmetric.Metric
	histogram pmetric.Histogram
}

func NewHistogramModel(logger *zap.Logger, warehouse, schema, table string) *HistogramModel {
	return &HistogramModel{
		logger:    logger,
		warehouse: warehouse,
		insertSQL: fmt.Sprintf(insertIntoHistogramMeticTableTemplate, schema, table),
	}
}

func (hm *HistogramModel) AddMetric(r pmetric.ResourceMetrics, s pmetric.ScopeMetrics, m pmetric.Metric) {
	hm.histograms = append(hm.histograms, &histogramData{
		resource:  r,
		scope:     s,
		metric:    m,
		histogram: m.Histogram(),
	})
}

func (hm *HistogramModel) BatchInsert(ctx context.Context, db database.Database) error {
	hm.logger.Debug("starting HistogramModel BatchInsert")
	if len(hm.histograms) == 0 {
		hm.logger.Debug("end HistogramModel BatchInsert: no histogram metrics to insert")
		return nil
	}

	histogramMaps := []map[string]any{}
	for _, h := range hm.histograms {
		for i := 0; i < h.histogram.DataPoints().Len(); i++ {
			dp := h.histogram.DataPoints().At(i)

			eAttributes, eTimestamps, eTraceIDs, eSpanIDs, eValues := utility.FlattenExemplars(dp.Exemplars(), hm.logger)

			histogramMaps = append(histogramMaps, map[string]any{
				"rSchema":        h.resource.SchemaUrl(),
				"rDroppedCount":  h.resource.Resource().DroppedAttributesCount(),
				"rAttributes":    utility.ConvertAttributesToString(h.resource.Resource().Attributes(), hm.logger),
				"sSchema":        h.scope.SchemaUrl(),
				"sName":          h.scope.Scope().Name(),
				"sVersion":       h.scope.Scope().Version(),
				"sDroppedCount":  h.scope.Scope().DroppedAttributesCount(),
				"sAttributes":    utility.ConvertAttributesToString(h.scope.Scope().Attributes(), hm.logger),
				"mName":          h.metric.Name(),
				"mDescription":   h.metric.Description(),
				"mUnit":          h.metric.Unit(),
				"attributes":     utility.ConvertAttributesToString(dp.Attributes(), hm.logger),
				"startTimestamp": dp.Timestamp().AsTime(),
				"timestamp":      dp.Timestamp().AsTime(),
				"count":          dp.Count(),
				"sum":            dp.Sum(),
				"flags":          dp.Flags(),
				"min":            dp.Min(),
				"max":            dp.Max(),
				"bucketCounts":   dp.BucketCounts().AsRaw(),
				"explicitBounds": dp.ExplicitBounds(),
				"eAttributes":    eAttributes,
				"eTimestamps":    eTimestamps,
				"eTraceIDs":      eTraceIDs,
				"eSpanIDs":       eSpanIDs,
				"eValues":        eValues,
			})
		}
	}

	hm.logger.Debug("HistogramModel calling utility.batchInsert")
	err := db.BatchInsert(ctx, histogramMaps, hm.warehouse, hm.insertSQL)
	if err != nil {
		return fmt.Errorf("failed to insert histogram metric data: %w", err)
	}
	hm.logger.Debug("end HistogramModel BatchInsert: successful insert")
	return nil
}
