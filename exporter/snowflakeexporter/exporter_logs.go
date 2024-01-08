package snowflakeexporter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/utility"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

const (
	createLogsTableSnowflakeSQL = `
	CREATE TABLE IF NOT EXISTS %s (
		"TraceId" BINARY(16),
		"SpanId" BINARY(16),
		"TraceFlags" BINARY(1),
		"SeverityText" VARCHAR,
		"SeverityNumber" NUMBER(2,0),
		"Body" VARIANT,
		"Resource" OBJECT,
		"InstrumentationScope" VARCHAR,
		"Attributes" OBJECT
	);
	`
)

type logsExporter struct {
	cfg    *Config
	logger *zap.Logger
	db     *sql.DB
}

func newLogsExporter(c *Config, params exporter.CreateSettings) (*logsExporter, error) {
	return &logsExporter{
		cfg:    c,
		logger: params.Logger,
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

	_, err = le.db.ExecContext(ctx, fmt.Sprintf(createLogsTableSnowflakeSQL, le.cfg.Logs.Table))
	if err != nil {
		return fmt.Errorf("failed to create table for logs: %w", err)
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
	return nil
}
