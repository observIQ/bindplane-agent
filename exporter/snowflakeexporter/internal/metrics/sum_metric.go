package metrics

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/utility"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

const (
	CreateSumMetricTableTemplate = `
	CREATE TABLE IF NOT EXISTS "%s"."%s_sum" (
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
		"IsMonotonic" BOOLEAN,
		"Attributes" VARCHAR,
		"StartTimestamp" TIMESTAMP_NTZ,
		"Timestamp" TIMESTAMP_NTZ,
		"Value" NUMBER,
		"Flags" INT,
		"ExemplarAttributes" VARCHAR,
		"ExemplarTimestamps" VARCHAR,
		"ExemplarTraceIDs" VARCHAR,
		"ExemplarSpanIDs" VARCHAR,
		"ExemplarValues" VARCHAR
	);`

	insertIntoSumMetricTableTemplate = `
	INSERT INTO "%s"."%s_sum" (
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
		"AggregationTemporality",
		"IsMonotonic",
		"Attributes",
		"StartTimestamp",
		"Timestamp",
		"Value",
		"Flags",
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
		:aggTemp,
		:monotonic,
		:attributes,
		:startTimestamp,
		:timestamp,
		:value,
		:flags,
		:eAttributes,
		:eTimestamps,
		:eTraceIDs,
		:eSpanIDs,
		:eValues
	);`
)

// SumModel implements the MetricModel for sum metrics
type SumModel struct {
	logger    *zap.Logger
	db        *sqlx.DB
	sums      []*sumData
	warehouse string
	insertSQL string
}

type sumData struct {
	resource pmetric.ResourceMetrics
	scope    pmetric.ScopeMetrics
	metric   pmetric.Metric
	sum      pmetric.Sum
}

// NewSumModel returns a new SumModel to be used for sending sum metrics to Snowflake
func NewSumModel(logger *zap.Logger, db *sqlx.DB, warehouse, schema, table string) *SumModel {
	return &SumModel{
		logger:    logger,
		db:        db,
		sums:      []*sumData{},
		warehouse: warehouse,
		insertSQL: fmt.Sprintf(insertIntoSumMetricTableTemplate, schema, table),
	}
}

// AddMetric adds a new sum metric to be inserted
func (sm *SumModel) AddMetric(r pmetric.ResourceMetrics, s pmetric.ScopeMetrics, m pmetric.Metric, d any) {
	sm.sums = append(sm.sums, &sumData{
		resource: r,
		scope:    s,
		metric:   m,
		sum:      d.(pmetric.Sum),
	})
}

// BatchInsert inserts all stored sum metrics
func (sm *SumModel) BatchInsert(ctx context.Context) error {
	sm.logger.Debug("starting SumModel BatchInsert")
	if len(sm.sums) == 0 {
		sm.logger.Debug("end SumModel BatchInsert: no sum metrics to insert")
		return nil
	}

	sumMaps := []map[string]any{}

	for _, s := range sm.sums {
		for i := 0; i < s.sum.DataPoints().Len(); i++ {
			dp := s.sum.DataPoints().At(i)

			var value any
			if dp.ValueType() == pmetric.NumberDataPointValueTypeInt {
				value = dp.IntValue()
			} else {
				value = dp.DoubleValue()
			}

			eAttributes, eTimestamps, eTraceIDs, eSpanIDs, eValues := utility.FlattenExemplars(dp.Exemplars(), *sm.logger)

			sumMaps = append(sumMaps, map[string]any{
				"rSchema":        s.resource.SchemaUrl(),
				"rDroppedCount":  s.resource.Resource().DroppedAttributesCount(),
				"rAttributes":    utility.ConvertAttributesToString(s.resource.Resource().Attributes(), sm.logger),
				"sSchema":        s.scope.SchemaUrl(),
				"sName":          s.scope.Scope().Name(),
				"sVersion":       s.scope.Scope().Version(),
				"sDroppedCount":  s.scope.Scope().DroppedAttributesCount(),
				"sAttributes":    utility.ConvertAttributesToString(s.scope.Scope().Attributes(), sm.logger),
				"mName":          s.metric.Name(),
				"mDescription":   s.metric.Description(),
				"mUnit":          s.metric.Unit(),
				"aggTemp":        s.sum.AggregationTemporality().String(),
				"monotonic":      s.sum.IsMonotonic(),
				"attributes":     utility.ConvertAttributesToString(dp.Attributes(), sm.logger),
				"startTimestamp": dp.StartTimestamp().AsTime(),
				"timestamp":      dp.Timestamp().AsTime(),
				"value":          value,
				"flags":          dp.Flags(),
				"eAttributes":    eAttributes,
				"eTimestamps":    eTimestamps,
				"eTraceIDs":      eTraceIDs,
				"eSpanIDs":       eSpanIDs,
				"eValues":        eValues,
			})
		}
	}

	sm.logger.Debug("SumModel calling utility.batchInsert")
	err := utility.BatchInsert(ctx, sm.db, sumMaps, sm.warehouse, sm.insertSQL)
	if err != nil {
		return fmt.Errorf("failed to insert sum metric data: %w", err)
	}
	sm.logger.Debug("end SumModel BatchInsert: successful insert")
	return nil
}
