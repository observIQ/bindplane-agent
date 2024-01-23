// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snowflakeexporter

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/snowflakedb/gosnowflake" // imports snowflake driver

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
		"ResourceSchemaURL" VARCHAR,
		"ResourceDroppedAttributesCount" INT,
		"ResourceAttributes" VARIANT,
		"ScopeSchemaURL" VARCHAR,
		"ScopeName" VARCHAR,
		"ScopeVersion" VARCHAR,
		"ScopeDroppedAttributesCount" INT,
		"ScopeAttributes" VARIANT,
		"TraceID" VARCHAR,
		"SpanID" VARCHAR,
		"TraceState" VARCHAR,
		"ParentSpanID" BINARY,
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
		Column1 AS "ResourceSchemaURL",
		Column2 AS "ResourceDroppedAttributesCount",
		PARSE_JSON(Column3) AS "ResourceAttributes",
		Column4 AS "ScopeSchemaURL",
		Column5 AS "ScopeName",
		Column6 AS "ScopeVersion",
		Column7 AS "ScopeDroppedAttributesCount",
		PARSE_JSON(Column8) AS "ScopeAttributes",
		Column9 AS "TraceID",
		Column10 AS "SpanID",
		Column11 AS "TraceState",
		Column12 AS "ParentSpanID",
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
		:rSchema,
		:rDroppedCount,
		:rAttributes,
		:sSchema,
		:sName,
		:sVersion,
		:sDroppedCount,
		:sAttributes,
		:traceID,
		:spanID,
		:traceState,
		:parentSpanID,
		:name,
		:kind,
		:startTime,
		:endTime,
		:droppedCount,
		:attributes,
		:statusMessage,
		:statusCode,
		:eventTimes,
		:eventNames,
		:eventAttributes,
		:eventDroppedCount,
		:linkTraceIDs,
		:linkSpanIDs,
		:linkTraceStates,
		:linkAttributes,
		:linkDroppedCount
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
	dsn := utility.CreateDSN(
		te.cfg.Username,
		te.cfg.Password,
		te.cfg.AccountIdentifier,
		te.cfg.Database,
		te.cfg.Traces.Schema,
	)
	db, err := utility.CreateDB(ctx, dsn)
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
		resourceSpan := td.ResourceSpans().At(i)

		for j := 0; j < resourceSpan.ScopeSpans().Len(); j++ {
			scopeSpan := resourceSpan.ScopeSpans().At(j)

			for k := 0; k < scopeSpan.Spans().Len(); k++ {
				span := scopeSpan.Spans().At(k)
				eTimes, eNames, eAttributes, eDroppedCount := flattenEvents(span.Events(), te.logger)
				lTraceIDs, lSpanIDs, lTraceStates, lAttributes, lDroppedCount := flattenLinks(span.Links(), te.logger)

				traceMaps = append(traceMaps, map[string]any{
					"rSchema":           resourceSpan.SchemaUrl(),
					"rDroppedCount":     resourceSpan.Resource().DroppedAttributesCount(),
					"rAttributes":       utility.ConvertAttributesToString(resourceSpan.Resource().Attributes(), te.logger),
					"sSchema":           scopeSpan.SchemaUrl(),
					"sName":             scopeSpan.Scope().Name(),
					"sVersion":          scopeSpan.Scope().Version(),
					"sDroppedCount":     scopeSpan.Scope().DroppedAttributesCount(),
					"sAttributes":       utility.ConvertAttributesToString(scopeSpan.Scope().Attributes(), te.logger),
					"traceID":           utility.TraceIDToHexOrEmptyString(span.TraceID()),
					"spanID":            utility.SpanIDToHexOrEmptyString(span.SpanID()),
					"traceState":        span.TraceState().AsRaw(),
					"parentSpanID":      utility.SpanIDToHexOrEmptyString(span.ParentSpanID()),
					"name":              span.Name(),
					"kind":              span.Kind().String(),
					"startTime":         span.StartTimestamp().AsTime(),
					"endTime":           span.EndTimestamp().AsTime(),
					"droppedCount":      span.DroppedAttributesCount(),
					"attributes":        utility.ConvertAttributesToString(span.Attributes(), te.logger),
					"statusMessage":     span.Status().Message(),
					"statusCode":        span.Status().Code().String(),
					"eventTimes":        eTimes,
					"eventNames":        eNames,
					"eventAttributes":   eAttributes,
					"eventDroppedCount": eDroppedCount,
					"linkTraceIDs":      lTraceIDs,
					"linkSpanIDs":       lSpanIDs,
					"linkTraceStates":   lTraceStates,
					"linkAttributes":    lAttributes,
					"linkDroppedCount":  lDroppedCount,
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

func flattenEvents(e ptrace.SpanEventSlice, l *zap.Logger) (pq.StringArray, pq.StringArray, pq.StringArray, pq.Int32Array) {
	times := pq.StringArray{}
	names := pq.StringArray{}
	attrs := pq.StringArray{}
	droppedCount := pq.Int32Array{}

	for i := 0; i < e.Len(); i++ {
		times = append(times, e.At(i).Timestamp().AsTime().String())
		names = append(names, e.At(i).Name())
		attrs = append(attrs, utility.ConvertAttributesToString(e.At(i).Attributes(), l))
		droppedCount = append(droppedCount, int32(e.At(i).DroppedAttributesCount()))
	}

	return times, names, attrs, droppedCount
}

func flattenLinks(li ptrace.SpanLinkSlice, l *zap.Logger) (pq.StringArray, pq.StringArray, pq.StringArray, pq.StringArray, pq.Int32Array) {
	traceIDs := pq.StringArray{}
	spanIDs := pq.StringArray{}
	traceStates := pq.StringArray{}
	attrs := pq.StringArray{}
	droppedCount := pq.Int32Array{}

	for i := 0; i < li.Len(); i++ {
		traceIDs = append(traceIDs, utility.TraceIDToHexOrEmptyString(li.At(i).TraceID()))
		spanIDs = append(spanIDs, utility.SpanIDToHexOrEmptyString(li.At(i).SpanID()))
		traceStates = append(traceStates, li.At(i).TraceState().AsRaw())
		attrs = append(attrs, utility.ConvertAttributesToString(li.At(i).Attributes(), l))
		droppedCount = append(droppedCount, int32(li.At(i).DroppedAttributesCount()))
	}

	return traceIDs, spanIDs, traceStates, attrs, droppedCount
}
