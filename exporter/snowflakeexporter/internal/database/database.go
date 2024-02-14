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

// Package database defines functions to be used by the Snowflake exporter for interacting with Snowflake
package database // import "github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/database"

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/snowflakedb/gosnowflake" // imports snowflake driver
)

// Database defines functions to use to interact with a database
//
//go:generate mockery --name Database --filename mock_database.go --structname MockDatabase
type Database interface {
	InitDatabaseConn(ctx context.Context, role, database, warehouse string) error
	CreateSchema(ctx context.Context, schema string) error
	CreateTable(ctx context.Context, table, template string) error
	BatchInsert(ctx context.Context, data []map[string]any, insertSQL string) error
	Close() error
}

// Snowflake implements the Database type
type Snowflake struct {
	db *sqlx.DB
}

// CreateSnowflakeDatabase calls Open() using driverName and the given dsn and then calls Ping()
func CreateSnowflakeDatabase(dsn string) (Database, error) {
	db, err := sqlx.Open("snowflake", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Snowflake{db: db}, nil
}

// InitDatabaseConn initializes the Snowflake connection by ensuring the correct role is used, database exists and is used by the connection, and that the warehouse will be used.
func (s *Snowflake) InitDatabaseConn(ctx context.Context, role, database, warehouse string) error {
	if role != "" {
		_, err := s.db.ExecContext(ctx, fmt.Sprintf(`USE ROLE "%s";`, role))
		if err != nil {
			return fmt.Errorf("failed to call 'USE ROLE \"%s\";': %w", role, err)
		}
	}

	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS "%s";`, database))
	if err != nil {
		return fmt.Errorf("failed to create new database: %w", err)
	}

	_, err = s.db.ExecContext(ctx, fmt.Sprintf(`USE DATABASE "%s";`, database))
	if err != nil {
		return fmt.Errorf("failed to call 'USE DATABASE \"%s\";': %w", database, err)
	}

	_, err = s.db.ExecContext(ctx, fmt.Sprintf(`USE WAREHOUSE "%s";`, warehouse))
	if err != nil {
		return fmt.Errorf("failed to call 'USE WAREHOUSE \"%s\";': %w", warehouse, err)
	}

	return nil
}

// CreateSchema ensures the given schema exists using the given *sql.DB
func (s *Snowflake) CreateSchema(ctx context.Context, schema string) error {
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, schema))
	if err != nil {
		return fmt.Errorf("failed to create schema '%s': %w", schema, err)
	}

	_, err = s.db.ExecContext(ctx, fmt.Sprintf(`USE SCHEMA "%s";`, schema))
	if err != nil {
		return fmt.Errorf("failed to call 'USE SCHEMA \"%s\";': %w", schema, err)
	}
	return nil
}

// CreateTable ensures the given table exists using the given database arguments
func (s *Snowflake) CreateTable(ctx context.Context, table, template string) error {
	_, err := s.db.ExecContext(ctx, fmt.Sprintf(template, table))
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// BatchInsert creates a new transaction using the given DB to insert the given data
func (s *Snowflake) BatchInsert(ctx context.Context, data []map[string]any, insertSQL string) error {
	// create TX
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.NamedExecContext(ctx, insertSQL, data)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}

	return tx.Commit()
}

// Close ensures the db is closed if it exists
func (s *Snowflake) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
