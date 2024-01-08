package snowflakeexporter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/utility"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

const (
	createTracesTableSnowflakeTemplate = `
	CREATE TABLE IF NOT EXISTS %s (
		"ResourceSchemaUrl" VARCHAR,
		"ResourceDroppedAttributesCount" INT,
		"ResourceAttributes" OBJECT,
		"ScopeSchemaUrl" VARCHAR,
		"ScopeName" VARCHAR,
		"ScopeVersion" VARCHAR,
		"ScopeDroppedAttributesCount" INT,
		"ScopeAttributes" OBJECT,
		"TraceId" BINARY,
		"SpanId" BINARY,
		"TraceState" VARCHAR,
		"ParentSpanId" BINARY,
		"Name" VARCHAR,
		"Kind" VARCHAR,
		"DroppedAttributesCount" INT,
		"Attributes" OBJECT,
		"StatusMessage" VARCHAR,
		"StatusCode" VARCHAR
	);
	`

	insertIntoTracesTableSnowflakeTemplate = `
	INSERT INTO %s (
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
		) VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
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
		return fmt.Errorf("failed to create new database for traces: %w", err)
	}
	te.db = db

	_, err = te.db.ExecContext(ctx, utility.RenderSQL(createTracesTableSnowflakeTemplate, te.cfg.Traces.Table))
	if err != nil {
		return fmt.Errorf("failed to create table for traces: %w", err)
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
	tx, err := te.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()
	if err := te.tracesTransaction(ctx, tx, td); err != nil {
		return err
	}
	return tx.Commit()
}

func (te *tracesExporter) tracesTransaction(ctx context.Context, tx *sql.Tx, td ptrace.Traces) error {
	stmt, err := tx.PrepareContext(ctx, utility.RenderSQL(insertIntoTracesTableSnowflakeTemplate, te.cfg.Traces.Table))
	if err != nil {
		return fmt.Errorf("failed to prepare transaction context: %w", err)
	}
	defer stmt.Close()

	for i := 0; i < td.ResourceSpans().Len(); i++ {
		resSpan := td.ResourceSpans().At(i)

		res := resSpan.Resource()
		resSchema := resSpan.SchemaUrl()
		resDroppedAttrsCount := res.DroppedAttributesCount()
		resAttrs := res.Attributes().AsRaw()
		for j := 0; j < resSpan.ScopeSpans().Len(); j++ {
			scopeSpan := resSpan.ScopeSpans().At(j)

			scopeSchema := scopeSpan.SchemaUrl()
			scopeName := scopeSpan.Scope().Name()
			scopeVersion := scopeSpan.Scope().Version()
			scopeDroppedAttrsCount := scopeSpan.Scope().DroppedAttributesCount()
			scopeAttrs := scopeSpan.Scope().Attributes().AsRaw()
			for k := 0; k < scopeSpan.Spans().Len(); k++ {
				span := scopeSpan.Spans().At(k)

				spanTraceID := utility.TraceIDToHexOrEmptyString(span.TraceID())
				spanSpanID := utility.SpanIDToHexOrEmptyString(span.SpanID())
				spanTraceState := span.TraceState().AsRaw()
				spanParentSpanID := utility.SpanIDToHexOrEmptyString(span.ParentSpanID())
				spanName := span.Name()
				spanKind := span.Kind().String()
				spanDroppedAttrsCount := span.DroppedAttributesCount()
				spanAttrs := span.Attributes().AsRaw()
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
