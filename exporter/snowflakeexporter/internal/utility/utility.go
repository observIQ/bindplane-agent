package utility

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
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
func CreateNewDB(ctx context.Context, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(DriverName, dsn)
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

	return db, nil
}

// CreateSchema ensures the given schema exists using the given *sql.DB
func CreateSchema(db *sqlx.DB, schema string) error {
	_, err := db.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, schema))
	if err != nil {
		return fmt.Errorf("failed to create schema '%s': %w", schema, err)
	}
	return nil
}

// CreateTable ensures the given table exists using the given database arguments
func CreateTable(ctx context.Context, db *sqlx.DB, database, schema, table, template string) error {
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

// BatchInsert creates a new transaction using the given DB to insert the given data
func BatchInsert(ctx context.Context, db *sqlx.DB, data *[]map[string]any, warehouse, insertSQL string) error {
	// create TX
	tx, err := db.BeginTxx(ctx, nil)
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
	_, err = tx.NamedExecContext(ctx, insertSQL, *data)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}

	return tx.Commit()
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
