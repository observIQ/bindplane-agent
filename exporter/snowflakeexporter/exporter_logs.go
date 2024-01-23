package snowflakeexporter

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/snowflakedb/gosnowflake"

	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/utility"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

const (
	createLogsTableSnowflakeTemplate = `
	CREATE TABLE IF NOT EXISTS "%s"."%s" (
		"ResourceSchemaURL" VARCHAR,
		"ResourceDroppedAttributesCount" INT,
		"ResourceAttributes" VARIANT,
		"ScopeSchemaUrl" VARCHAR,
		"ScopeName" VARCHAR,
		"ScopeVersion" VARCHAR,
		"ScopeDroppedAttributesCount" INT,
		"ScopeAttributes" VARIANT,
		"Timestamp" TIMESTAMP_NTZ,
		"ObservedTimestamp" TIMESTAMP_NTZ,
		"SeverityNumber" VARCHAR,
		"SeverityText" VARCHAR,
		"Body" VARCHAR,
		"Attributes" VARIANT,
		"DroppedAttributesCount" INT,
		"Flags" INT,
		"TraceID" VARCHAR,
		"SpanID" VARCHAR
	);`

	insertIntoLogsTableSnowflakeTemplate = `
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
		Column9 AS "Timestamp",
		Column10 AS "ObservedTimestamp",
		Column11 AS "SeverityNumber",
		Column12 AS "SeverityText",
		Column13 AS "Body",
		PARSE_JSON(Column14) AS "Attributes",
		Column15 AS "DroppedAttributesCount",
		Column16 AS "Flags",
		Column17 AS "TraceID",
		Column18 AS "SpanID"
	FROM VALUES (
		:rSchema,
		:rDroppedCount,
		:rAttributes,
		:sSchema,
		:sName,
		:sVersion,
		:sDroppedCount,
		:sAttributes,
		:timestamp,
		:observedTimestamp,
		:severityNumber,
		:severityText,
		:body
		:attributes,
		:droppedCount,
		:flags,
		:traceID,
		:spanID
	);`
)

type logsExporter struct {
	cfg       *Config
	logger    *zap.Logger
	db        *sqlx.DB
	insertSQL string
}

func newLogsExporter(c *Config, params exporter.CreateSettings) (*logsExporter, error) {
	return &logsExporter{
		cfg:       c,
		logger:    params.Logger,
		insertSQL: fmt.Sprintf(insertIntoLogsTableSnowflakeTemplate, c.Logs.Schema, c.Logs.Table),
	}, nil
}

func (le *logsExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

func (le *logsExporter) start(ctx context.Context, _ component.Host) error {
	dsn := utility.BuildDSN(
		le.cfg.Username,
		le.cfg.Password,
		le.cfg.AccountIdentifier,
		le.cfg.Database,
		le.cfg.Logs.Schema,
	)
	db, err := utility.CreateNewDB(ctx, dsn)
	if err != nil {
		return fmt.Errorf("failed to create new database for logs: %w", err)
	}
	le.db = db

	err = utility.CreateTable(ctx, le.db, le.cfg.Database, le.cfg.Logs.Schema, le.cfg.Logs.Table, createLogsTableSnowflakeTemplate)
	if err != nil {
		le.logger.Debug("CreateTable failed for logs", zap.String("database", le.cfg.Database), zap.String("schema", le.cfg.Logs.Schema), zap.String("table", le.cfg.Logs.Table))
		return fmt.Errorf("failed to create logs table: %w", err)
	}

	return nil
}

func (le *logsExporter) shutdown(_ context.Context) error {
	if le.db != nil {
		return le.db.Close()
	}
	return nil
}

// entry function
func (le *logsExporter) logsDataPusher(ctx context.Context, ld plog.Logs) error {
	le.logger.Debug("begin logsDataPusher")
	logMaps := []map[string]any{}

	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resourceLog := ld.ResourceLogs().At(i)

		for j := 0; j < resourceLog.ScopeLogs().Len(); j++ {
			scopeLog := resourceLog.ScopeLogs().At(j)

			for k := 0; k < scopeLog.LogRecords().Len(); k++ {
				log := scopeLog.LogRecords().At(k)

				logMaps = append(logMaps, map[string]any{
					"rSchema":           resourceLog.SchemaUrl(),
					"rDroppedCount":     resourceLog.Resource().DroppedAttributesCount(),
					"rAttributes":       utility.ConvertAttributesToString(resourceLog.Resource().Attributes(), le.logger),
					"sSchema":           scopeLog.SchemaUrl(),
					"sName":             scopeLog.Scope().Name(),
					"sVersion":          scopeLog.Scope().Version(),
					"sDroppedCount":     scopeLog.Scope().DroppedAttributesCount(),
					"sAttributes":       utility.ConvertAttributesToString(scopeLog.Scope().Attributes(), le.logger),
					"timestamp":         log.Timestamp().AsTime(),
					"observedTimestamp": log.ObservedTimestamp().AsTime(),
					"severityNumber":    log.SeverityNumber().String(),
					"severityText":      log.SeverityText(),
					"body":              log.Body().AsString(),
					"attributes":        utility.ConvertAttributesToString(log.Attributes(), le.logger),
					"droppedCount":      log.DroppedAttributesCount(),
					"flags":             int32(log.Flags()),
					"traceID":           utility.TraceIDToHexOrEmptyString(log.TraceID()),
					"spanID":            utility.SpanIDToHexOrEmptyString(log.SpanID()),
				})
			}
		}
	}

	le.logger.Debug("insert sql", zap.String("logs", le.insertSQL))
	err := utility.BatchInsert(ctx, le.db, &logMaps, le.cfg.Warehouse, le.insertSQL)
	if err != nil {
		return fmt.Errorf("failed to insert log data: %w", err)
	}
	le.logger.Debug("end logsDataPusher")
	return nil
}
