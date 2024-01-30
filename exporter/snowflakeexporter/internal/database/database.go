package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/snowflakedb/gosnowflake" // imports snowflake driver
)

//go:generate mockery --name Database --filename mock_database.go --structname MockDatabase
type Database interface {
	CreateSchema(ctx context.Context, schema string) error
	CreateTable(ctx context.Context, database, schema, table, template string) error
	BatchInsert(ctx context.Context, data []map[string]any, warehouse, insertSQL string) error
	Close() error
}

type Snowflake struct {
	db *sqlx.DB
}

// CreateSnowflakeDatabase calls Open() using driverName and the given dsn and then calls Ping()
func CreateSnowflakeDatabase(ctx context.Context, dsn string) (Database, error) {
	db, err := sqlx.Open("snowflake", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if ctx == nil {
		if err = db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}
	} else {
		if err = db.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}
	}

	return &Snowflake{db: db}, nil
}

// CreateSchema ensures the given schema exists using the given *sql.DB
func (s *Snowflake) CreateSchema(ctx context.Context, schema string) error {
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, schema))
	if err != nil {
		return fmt.Errorf("failed to create schema '%s': %w", schema, err)
	}
	return nil
}

// CreateTable ensures the given table exists using the given database arguments
func (s *Snowflake) CreateTable(ctx context.Context, database, schema, table, template string) error {
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`USE SCHEMA "%s"."%s";`, database, schema))
	if err != nil {
		return fmt.Errorf("failed to call 'USE SCHEMA': %w", err)
	}

	_, err = s.db.ExecContext(ctx, fmt.Sprintf(template, schema, table))
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// BatchInsert creates a new transaction using the given DB to insert the given data
func (s *Snowflake) BatchInsert(ctx context.Context, data []map[string]any, warehouse, insertSQL string) error {
	// create TX
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	// set warehouse
	_, err = tx.ExecContext(ctx, fmt.Sprintf(`USE WAREHOUSE "%s";`, warehouse))
	if err != nil {
		return fmt.Errorf("failed to set warehouse as '%s' for transaction: %w", warehouse, err)
	}

	// execute insert
	_, err = tx.NamedExecContext(ctx, insertSQL, data)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}

	return tx.Commit()
}

func (s *Snowflake) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
