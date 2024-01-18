package snowflakeexporter

import (
	"context"
	"database/sql"
	"fmt"

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
		"TraceId" BINARY,
		"SpanId" BINARY,
		"TraceState" VARCHAR,
		"ParentSpanId" BINARY,
		"Name" VARCHAR,
		"Kind" VARCHAR,
		"DroppedAttributesCount" INT,
		"Attributes" VARIANT,
		"StatusMessage" VARCHAR,
		"StatusCode" VARCHAR
	);
	`

	insertIntoTracesTableSnowflakeTemplate = `
	INSERT INTO "%s"."%s" (
		"ResourceSchemaUrl",
		"ResourceDroppedAttributesCount",
		"ResourceAttributes",
		"ScopeSchemaUrl",
		"ScopeName",
		"ScopeVersion",
		"ScopeDroppedAttributesCount",
		"ScopeAttributes",
		"TraceId",
		"SpanId",
		"TraceState",
		"ParentSpanId",
		"Name",
		"Kind",
		"DroppedAttributesCount",
		"Attributes",
		"StatusMessage",
		"StatusCode"
		) SELECT
			?,
			?,
			PARSE_JSON(?),
			?,
			?,
			?,
			?,
			PARSE_JSON(?),
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			PARSE_JSON(?),
			?,
			?
		;
	`
)

type tracesExporter struct {
	cfg    *Config
	logger *zap.Logger
	db     *sql.DB
}

func newTracesExporter(c *Config, params exporter.CreateSettings) (*tracesExporter, error) {
	return &tracesExporter{
		cfg:    c,
		logger: params.Logger,
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
	te.logger.Debug("begin transaction")
	tx, err := te.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, fmt.Sprintf(`USE WAREHOUSE "%s";`, te.cfg.Warehouse))
	if err != nil {
		return fmt.Errorf("failed to call 'USE WAREHOUSE': %w", err)
	}

	if err := te.tracesTransaction(ctx, tx, td); err != nil {
		te.logger.Debug("failed transaction", zap.Error(err))
		return err
	}
	te.logger.Debug("successful transaction")
	return tx.Commit()
}

func (te *tracesExporter) tracesTransaction(ctx context.Context, tx *sql.Tx, td ptrace.Traces) error {
	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf(insertIntoTracesTableSnowflakeTemplate, te.cfg.Traces.Schema, te.cfg.Traces.Table))
	if err != nil {
		return fmt.Errorf("failed to prepare transaction context: %w", err)
	}
	defer stmt.Close()

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
				spanDroppedAttrsCount := span.DroppedAttributesCount()
				spanAttrs := utility.ConvertAttributesToString(span.Attributes(), te.logger)
				spanStatusMessage := span.Status().Message()
				spanStatusCode := span.Status().Code().String()

				_, err := stmt.ExecContext(ctx,
					resSchema,
					resDroppedAttrsCount,
					resAttrs,
					scopeSchema,
					scopeName,
					scopeVersion,
					scopeDroppedAttrsCount,
					scopeAttrs,
					spanTraceID,
					spanSpanID,
					spanTraceState,
					spanParentSpanID,
					spanName,
					spanKind,
					spanDroppedAttrsCount,
					spanAttrs,
					spanStatusMessage,
					spanStatusCode,
				)
				if err != nil {
					return fmt.Errorf("failed to execute statement: %w", err)
				}
			}
		}
	}
	return nil
}
