package utility

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/url"

	"go.opentelemetry.io/collector/pdata/pcommon"
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
		dsn = fmt.Sprintf("%s:%s@%s/%s", usernameEsc, passwordEsc, accountIDEsc, databaseEsc)
	} else {
		schemaEsc := url.QueryEscape(schema)
		dsn = fmt.Sprintf("%s:%s@%s/%s/%s", usernameEsc, passwordEsc, accountIDEsc, databaseEsc, schemaEsc)
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
	_, err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s;", schema))
	if err != nil {
		return fmt.Errorf("failed to create schema '%s': %w", schema, err)
	}
	return nil
}

func RenderSQL(template string, s string) string {
	return fmt.Sprintf(template, s)
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
