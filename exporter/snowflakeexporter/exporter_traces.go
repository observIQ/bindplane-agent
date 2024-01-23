package snowflakeexporter

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/snowflakedb/gosnowflake"

	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/utility"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

const (
	createTracesTableSnowflakeTemplate = `
	CREATE TABLE IF NOT EXISTS "%s"."%s" (
		"ResourceSchemaUrl" VARCHAR,
		"ResourceDroppedAttributesCount" INT,
		"ResourceAttributes" VARIANT,
		"ScopeSchemaUrl" VARCHAR,
		"ScopeName" VARCHAR,
		"ScopeVersion" VARCHAR,
		"ScopeDroppedAttributesCount" INT,
		"ScopeAttributes" VARIANT,
		"TraceID" VARCHAR,
		"SpanID" VARCHAR,
		"TraceState" VARCHAR,
		"ParentSpanId" BINARY,
		"Name" VARCHAR,
		"Kind" VARCHAR,
		"StartTime" TIMESTAMP_NTZ,
		"EndTime" TIMESTAMP_NTZ,
		"DroppedAttributesCount" INT,
		"Attributes" VARIANT,
		"StatusMessage" VARCHAR,
		"StatusCode" VARCHAR,
		"EventTimes" ARRAY,
		"EventNames" ARRAY,
		"EventAttributes" ARRAY,
		"EventDroppedAttributesCount" ARRAY,
		"LinkTraceID" ARRAY,
		"LinkSpanID" ARRAY,
		"LinkTraceState" ARRAY,
		"LinkDroppedAttributesCount" ARRAY,
		"LinkAttributes" ARRAY
	);`

	insertIntoTracesTableSnowflakeTemplate = `
	INSERT INTO "%s"."%s"
	SELECT
		Column1 AS "ResourceSchemaUrl",
		Column2 AS "ResourceDroppedAttributesCount",
		PARSE_JSON(Column3) AS "ResourceAttributes",
		Column4 AS "ScopeSchemaUrl",
		Column5 AS "ScopeName",
		Column6 AS "ScopeVersion",
		Column7 AS "ScopeDroppedAttributesCount",
		PARSE_JSON(Column8) AS "ScopeAttributes",
		Column9 AS "TraceId",
		Column10 AS "SpanId",
		Column11 AS "TraceState",
		Column12 AS "ParentSpanId",
		Column13 AS "Name",
		Column14 AS "Kind",
		Column15 AS "StartTime",
		Column16 AS "EndTime",
		Column17 AS "DroppedAttributesCount",
		PARSE_JSON(Column18) AS "Attributes",
		Column19 AS "StatusMessage",
		Column20 AS "StatusCode",
		ARRAY_CONSTRUCT(Column21) AS "EventTimes",
		ARRAY_CONSTRUCT(Column22) AS "EventNames",
		ARRAY_CONSTRUCT(Column23) AS "EventAttributes",
		ARRAY_CONSTRUCT(Column24) AS "EventDroppedAttributesCount",
		ARRAY_CONSTRUCT(Column25) AS "LinkTraceID",
		ARRAY_CONSTRUCT(Column26) AS "LinkSpanID",
		ARRAY_CONSTRUCT(Column27) AS "LinkTraceState",
		ARRAY_CONSTRUCT(Column28) AS "LinkDroppedAttributesCount",
		ARRAY_CONSTRUCT(Column29) AS "LinkAttributes"
	FROM VALUES (
		:resSchema,
		:resDroppedAttrsCount,
		:resAttrs,
		:scopeSchema,
		:scopeName,
		:scopeVersion,
		:scopeDroppedAttrsCount,
		:scopeAttrs,
		:spanTraceID,
		:spanSpanID,
		:spanTraceState,
		:spanParentSpanID,
		:spanName,
		:spanKind,
		:spanStartTime,
		:spanEndTime,
		:spanDroppedAttrsCount,
		:spanAttrs,
		:spanStatusMessage,
		:spanStatusCode,
		:spanEventTimes,
		:spanEventNames,
		:spanEventAttrs,
		:spanEventDroppedAttrs,
		:spanLinkTraceIDs,
		:spanLinkSpanIDs,
		:spanLinkTraceStates,
		:spanLinkAttrs,
		:spanLinkDroppedAttrs
	);`
)

type tracesExporter struct {
	cfg       *Config
	logger    *zap.Logger
	db        *sqlx.DB
	insertSQL string
}

func newTracesExporter(c *Config, params exporter.CreateSettings) (*tracesExporter, error) {
	return &tracesExporter{
		cfg:       c,
		logger:    params.Logger,
		insertSQL: fmt.Sprintf(insertIntoTracesTableSnowflakeTemplate, c.Traces.Schema, c.Traces.Table),
	}, nil
}

func (te *tracesExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (te *tracesExporter) start(ctx context.Context, _ component.Host) error {
	dsn := utility.BuildDSN(
		te.cfg.Username,
		te.cfg.Password,
		te.cfg.AccountIdentifier,
		te.cfg.Database,
		te.cfg.Traces.Schema,
	)
	db, err := utility.CreateNewDB(ctx, dsn)
	if err != nil {
		te.logger.Debug("CreateNewDB failed for traces", zap.String("dsn", dsn))
		return fmt.Errorf("failed to create new database for traces: %w", err)
	}
	te.db = db

	err = utility.CreateTable(ctx, te.db, te.cfg.Database, te.cfg.Traces.Schema, te.cfg.Traces.Table, createTracesTableSnowflakeTemplate)
	if err != nil {
		te.logger.Debug("CreateTable failed for traces", zap.String("database", te.cfg.Database), zap.String("schema", te.cfg.Traces.Schema), zap.String("table", te.cfg.Traces.Table))
		return fmt.Errorf("failed to create traces table: %w", err)
	}

	return nil
}

func (te *tracesExporter) shutdown(_ context.Context) error {
	if te.db != nil {
		return te.db.Close()
	}
	return nil
}

// entry function
func (te *tracesExporter) tracesDataPusher(ctx context.Context, td ptrace.Traces) error {
	te.logger.Debug("begin tracesDataPusher")
	traceMaps := []map[string]any{}
	for i := 0; i < td.ResourceSpans().Len(); i++ {
		resSpan := td.ResourceSpans().At(i)

		res := resSpan.Resource()
		resSchema := resSpan.SchemaUrl()
		resDroppedAttrsCount := res.DroppedAttributesCount()
		resAttrs := utility.ConvertAttributesToString(res.Attributes(), te.logger)
		for j := 0; j < resSpan.ScopeSpans().Len(); j++ {
			scopeSpan := resSpan.ScopeSpans().At(j)

			scopeSchema := scopeSpan.SchemaUrl()
			scopeName := scopeSpan.Scope().Name()
			scopeVersion := scopeSpan.Scope().Version()
			scopeDroppedAttrsCount := scopeSpan.Scope().DroppedAttributesCount()
			scopeAttrs := utility.ConvertAttributesToString(scopeSpan.Scope().Attributes(), te.logger)
			for k := 0; k < scopeSpan.Spans().Len(); k++ {
				span := scopeSpan.Spans().At(k)

				spanTraceID := utility.TraceIDToHexOrEmptyString(span.TraceID())
				spanSpanID := utility.SpanIDToHexOrEmptyString(span.SpanID())
				spanTraceState := span.TraceState().AsRaw()
				spanParentSpanID := utility.SpanIDToHexOrEmptyString(span.ParentSpanID())
				spanName := span.Name()
				spanKind := span.Kind().String()
				spanStartTime := span.StartTimestamp().AsTime()
				spanEndTime := span.EndTimestamp().AsTime()
				spanDroppedAttrsCount := span.DroppedAttributesCount()
				spanAttrs := utility.ConvertAttributesToString(span.Attributes(), te.logger)
				spanStatusMessage := span.Status().Message()
				spanStatusCode := span.Status().Code().String()
				eventTimes, eventNames, eventAttrs, eventDroppedAttrsCount := te.flattenEvents(span.Events())
				linkTraceIDs, linkSpanIDs, linkTraceStates, linkAttrs, linkDroppedAttrsCount := te.flattenLinks(span.Links())

				traceMaps = append(traceMaps, map[string]any{
					"resSchema":              resSchema,
					"resDroppedAttrsCount":   resDroppedAttrsCount,
					"resAttrs":               resAttrs,
					"scopeSchema":            scopeSchema,
					"scopeName":              scopeName,
					"scopeVersion":           scopeVersion,
					"scopeDroppedAttrsCount": scopeDroppedAttrsCount,
					"scopeAttrs":             scopeAttrs,
					"spanTraceID":            spanTraceID,
					"spanSpanID":             spanSpanID,
					"spanTraceState":         spanTraceState,
					"spanParentSpanID":       spanParentSpanID,
					"spanName":               spanName,
					"spanKind":               spanKind,
					"spanStartTime":          spanStartTime,
					"spanEndTime":            spanEndTime,
					"spanDroppedAttrsCount":  spanDroppedAttrsCount,
					"spanAttrs":              spanAttrs,
					"spanStatusMessage":      spanStatusMessage,
					"spanStatusCode":         spanStatusCode,
					"spanEventTimes":         eventTimes,
					"spanEventNames":         eventNames,
					"spanEventAttrs":         eventAttrs,
					"spanEventDroppedAttrs":  eventDroppedAttrsCount,
					"spanLinkTraceIDs":       linkTraceIDs,
					"spanLinkSpanIDs":        linkSpanIDs,
					"spanLinkTraceStates":    linkTraceStates,
					"spanLinkAttrs":          linkAttrs,
					"spanLinkDroppedAttrs":   linkDroppedAttrsCount,
				})
			}
		}
	}

	err := utility.BatchInsert(ctx, te.db, &traceMaps, te.cfg.Warehouse, te.insertSQL)
	if err != nil {
		return fmt.Errorf("failed to insert trace data: %w", err)
	}
	te.logger.Debug("end tracesDataPusher")
	return nil
}

func (te *tracesExporter) flattenEvents(e ptrace.SpanEventSlice) (pq.StringArray, pq.StringArray, pq.StringArray, pq.Int32Array) {
	times := pq.StringArray{}
	names := pq.StringArray{}
	attrs := pq.StringArray{}
	droppedCount := pq.Int32Array{}

	for i := 0; i < e.Len(); i++ {
		times = append(times, e.At(i).Timestamp().AsTime().String())
		names = append(names, e.At(i).Name())
		attrs = append(attrs, utility.ConvertAttributesToString(e.At(i).Attributes(), te.logger))
		droppedCount = append(droppedCount, int32(e.At(i).DroppedAttributesCount()))
	}

	return times, names, attrs, droppedCount
}

func (te *tracesExporter) flattenLinks(l ptrace.SpanLinkSlice) (pq.StringArray, pq.StringArray, pq.StringArray, pq.StringArray, pq.Int32Array) {
	traceIDs := pq.StringArray{}
	spanIDs := pq.StringArray{}
	traceStates := pq.StringArray{}
	attrs := pq.StringArray{}
	droppedCount := pq.Int32Array{}

	for i := 0; i < l.Len(); i++ {
		traceIDs = append(traceIDs, utility.TraceIDToHexOrEmptyString(l.At(i).TraceID()))
		spanIDs = append(spanIDs, utility.SpanIDToHexOrEmptyString(l.At(i).SpanID()))
		traceStates = append(traceStates, l.At(i).TraceState().AsRaw())
		attrs = append(attrs, utility.ConvertAttributesToString(l.At(i).Attributes(), te.logger))
		droppedCount = append(droppedCount, int32(l.At(i).DroppedAttributesCount()))
	}

	return traceIDs, spanIDs, traceStates, attrs, droppedCount
}
