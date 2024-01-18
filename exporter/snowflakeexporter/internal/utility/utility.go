package utility

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"

	_ "github.com/snowflakedb/gosnowflake"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.uber.org/zap"
)

// allows use of mocking by changing
var DriverName = "snowflake"

// BuildDSN creates a DSN for connecting to Snowflake with the given config
// TODO add functionality for additional query params
func BuildDSN(username, password, accountID, database, schema string) string {
	var dsn string
	usernameEsc := url.QueryEscape(username)
	passwordEsc := url.QueryEscape(password)
	accountIDEsc := url.QueryEscape(accountID)
	databaseEsc := url.QueryEscape(database)
	if schema != "" {
		dsn = fmt.Sprintf(`%s:%s@%s/"%s"`, usernameEsc, passwordEsc, accountIDEsc, databaseEsc)
	} else {
		schemaEsc := url.QueryEscape(schema)
		dsn = fmt.Sprintf(`%s:%s@%s/"%s"/"%s"`, usernameEsc, passwordEsc, accountIDEsc, databaseEsc, schemaEsc)
	}
	return dsn
}

// CreateNewDB calls Open() using driverName and the given dsn and then calls Ping()
func CreateNewDB(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open(DriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

// CreateSchema ensures the given schema exists using the given *sql.DB
func CreateSchema(db *sql.DB, schema string) error {
	_, err := db.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, schema))
	if err != nil {
		return fmt.Errorf("failed to create schema '%s': %w", schema, err)
	}
	return nil
}

// CreateTable ensures teh given table exists using the given database arguments
func CreateTable(ctx context.Context, db *sql.DB, database, schema, table, template string) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf(`USE SCHEMA "%s"."%s";`, database, schema))
	if err != nil {
		return fmt.Errorf("failed to call 'USE SCHEMA': %w", err)
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf(template, schema, table))
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func TraceIDToHexOrEmptyString(id pcommon.TraceID) string {
	if id.IsEmpty() {
		return ""
	}
	return hex.EncodeToString(id[:])
}

func SpanIDToHexOrEmptyString(id pcommon.SpanID) string {
	if id.IsEmpty() {
		return ""
	}
	return hex.EncodeToString(id[:])
}

// ConvertAttributesToString converts the pcommon.Map into a JSON string representation
// this is due to a bug/lacking feature with the snowflake driver that prevents maps from being inserted into VARIANT & OBJECT columns
// github issue: https://github.com/snowflakedb/gosnowflake/issues/217
func ConvertAttributesToString(attributes pcommon.Map, logger *zap.Logger) string {
	m := make(map[string]string, attributes.Len())
	attributes.Range(func(k string, v pcommon.Value) bool {
		m[k] = v.AsString()
		return true
	})

	j, err := json.Marshal(m)
	if err != nil {
		logger.Warn("failed to marshal attribute map", zap.Error(err), zap.Any("map", m))
		return ""
	}
	return string(j)
}
