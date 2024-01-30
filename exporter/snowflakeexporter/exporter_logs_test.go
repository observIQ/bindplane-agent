package snowflakeexporter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/database"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/database/mocks"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.opentelemetry.io/collector/pdata/plog"
)

func TestNewLogsExporter(t *testing.T) {
	c := &Config{
		AccountIdentifier: "id",
		Username:          "user",
		Password:          "pass",
		Database:          "db",
		Logs: &TelemetryConfig{
			Schema: "schema",
			Table:  "table",
		},
	}

	testCases := []struct {
		desc        string
		ctx         context.Context
		newDatabase func(ctx context.Context, dsn string) (database.Database, error)
		expectedErr error
	}{
		{
			desc: "Simple pass",
			ctx:  context.Background(),
			newDatabase: func(_ context.Context, _ string) (database.Database, error) {
				return mocks.NewMockDatabase(t), nil
			},
		},
		{
			desc: "Fail newDatabase",
			ctx:  context.Background(),
			newDatabase: func(_ context.Context, _ string) (database.Database, error) {
				return nil, fmt.Errorf("fail")
			},
			expectedErr: fmt.Errorf("failed to create new database connection for logs: fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			exp, err := newLogsExporter(
				tc.ctx,
				c,
				exportertest.NewNopCreateSettings(),
				tc.newDatabase,
			)

			if tc.expectedErr == nil {
				require.NoError(t, err)
				require.NotNil(t, exp)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
				require.Nil(t, exp)
			}
		})
	}
}

func TestCapabilities(t *testing.T) {
	e := &logsExporter{}
	c := e.Capabilities()
	require.False(t, c.MutatesData)
}

func TestStart(t *testing.T) {
	c := &Config{
		AccountIdentifier: "id",
		Username:          "user",
		Password:          "pass",
		Database:          "db",
		Logs: &TelemetryConfig{
			Schema: "schema",
			Table:  "table",
		},
	}

	testCases := []struct {
		desc        string
		ctx         context.Context
		mockGen     func(t *testing.T, ctx context.Context, c *Config) *mocks.MockDatabase
		expectedErr error
	}{
		{
			desc: "Simple pass",
			ctx:  context.Background(),
			mockGen: func(t *testing.T, ctx context.Context, c *Config) *mocks.MockDatabase {
				m := mocks.NewMockDatabase(t)
				m.On("CreateSchema", ctx, c.Logs.Schema).Return(nil)
				m.On("CreateTable", ctx, c.Database, c.Logs.Schema, c.Logs.Table, createLogsTableSnowflakeTemplate).Return(nil)
				return m
			},
		},
		{
			desc:        "Fail CreateSchema",
			ctx:         context.Background(),
			expectedErr: fmt.Errorf("failed to create logs schema: fail"),
			mockGen: func(t *testing.T, ctx context.Context, c *Config) *mocks.MockDatabase {
				m := mocks.NewMockDatabase(t)
				m.On("CreateSchema", ctx, c.Logs.Schema).Return(fmt.Errorf("fail"))
				return m
			},
		},
		{
			desc:        "Fail CreateTable",
			ctx:         context.Background(),
			expectedErr: fmt.Errorf("failed to create logs table: fail"),
			mockGen: func(t *testing.T, ctx context.Context, c *Config) *mocks.MockDatabase {
				m := mocks.NewMockDatabase(t)
				m.On("CreateSchema", ctx, c.Logs.Schema).Return(nil)
				m.On("CreateTable", ctx, c.Database, c.Logs.Schema, c.Logs.Table, createLogsTableSnowflakeTemplate).Return(fmt.Errorf("fail"))
				return m
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			logsExp, err := newLogsExporter(
				tc.ctx,
				c,
				exportertest.NewNopCreateSettings(),
				func(ctx context.Context, dsn string) (database.Database, error) { return nil, nil },
			)
			require.NoError(t, err)
			logsExp.db = tc.mockGen(t, tc.ctx, c)

			err = logsExp.start(tc.ctx, nil)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}

func TestShutdown(t *testing.T) {
	// no db
	ctx := context.Background()
	e := &logsExporter{}
	require.NoError(t, e.shutdown(ctx))

	// db, no error & error
	mock := mocks.NewMockDatabase(t)
	mock.On("Close").Return(nil).Once()
	mock.On("Close").Return(fmt.Errorf("fail")).Once()
	e.db = mock

	require.NoError(t, e.shutdown(ctx))
	require.ErrorContains(t, e.shutdown(ctx), "fail")
}

func TestLogsDataPusher(t *testing.T) {
	c := &Config{
		AccountIdentifier: "id",
		Username:          "user",
		Password:          "pass",
		Database:          "db",
		Logs: &TelemetryConfig{
			Schema: "schema",
			Table:  "table",
		},
	}

	testCases := []struct {
		desc        string
		ctx         context.Context
		logGen      func(t *testing.T) plog.Logs
		mapGen      func(t *testing.T) []map[string]any
		mockGen     func(t *testing.T, ctx context.Context, warehouse, sql string, data []map[string]any) *mocks.MockDatabase
		expectedErr error
	}{
		{
			desc:   "Simple pass",
			ctx:    context.Background(),
			logGen: generateLogData1,
			mapGen: generateLogMaps1,
			mockGen: func(t *testing.T, ctx context.Context, warehouse, sql string, data []map[string]any) *mocks.MockDatabase {
				mock := mocks.NewMockDatabase(t)
				mock.On("BatchInsert", ctx, data, c.Warehouse, sql).Return(nil)
				return mock
			},
		},
		{
			desc:   "Simple fail",
			ctx:    context.Background(),
			logGen: generateLogData1,
			mapGen: generateLogMaps1,
			mockGen: func(t *testing.T, ctx context.Context, warehouse, sql string, data []map[string]any) *mocks.MockDatabase {
				mock := mocks.NewMockDatabase(t)
				mock.On("BatchInsert", ctx, data, c.Warehouse, sql).Return(fmt.Errorf("fail"))
				return mock
			},
			expectedErr: fmt.Errorf("failed to insert log data: fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			logsExp, err := newLogsExporter(
				tc.ctx,
				c,
				exportertest.NewNopCreateSettings(),
				func(ctx context.Context, dsn string) (database.Database, error) { return nil, nil },
			)
			require.NoError(t, err)
			logsExp.db = tc.mockGen(t, tc.ctx, c.Warehouse, logsExp.insertSQL, tc.mapGen(t))

			err = logsExp.logsDataPusher(tc.ctx, tc.logGen(t))
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}

func generateLogData1(t *testing.T) plog.Logs {
	logs := plog.NewLogs()
	rLogs := logs.ResourceLogs().AppendEmpty()
	rLogs.SetSchemaUrl("resource_schema_url")
	sLogs := rLogs.ScopeLogs().AppendEmpty()
	sLogs.SetSchemaUrl("scope_schema_url")
	for i := 0; i < 3; i++ {
		lr := sLogs.LogRecords().AppendEmpty()
		lr.Body().SetStr(fmt.Sprintf("log_body_%d", i))
	}
	return logs
}

func generateLogMaps1(t *testing.T) []map[string]any {
	return []map[string]any{
		{
			"rSchema":           "resource_schema_url",
			"rDroppedCount":     uint32(0),
			"rAttributes":       "{}",
			"sSchema":           "scope_schema_url",
			"sName":             "",
			"sVersion":          "",
			"sDroppedCount":     uint32(0),
			"sAttributes":       "{}",
			"timestamp":         time.Unix(0, int64(0)).UTC(),
			"observedTimestamp": time.Unix(0, int64(0)).UTC(),
			"severityNumber":    "Unspecified",
			"severityText":      "",
			"body":              "log_body_0",
			"attributes":        "{}",
			"droppedCount":      uint32(0),
			"flags":             plog.LogRecordFlags(0),
			"traceID":           "",
			"spanID":            "",
		},
		{
			"rSchema":           "resource_schema_url",
			"rDroppedCount":     uint32(0),
			"rAttributes":       "{}",
			"sSchema":           "scope_schema_url",
			"sName":             "",
			"sVersion":          "",
			"sDroppedCount":     uint32(0),
			"sAttributes":       "{}",
			"timestamp":         time.Unix(0, int64(0)).UTC(),
			"observedTimestamp": time.Unix(0, int64(0)).UTC(),
			"severityNumber":    "Unspecified",
			"severityText":      "",
			"body":              "log_body_1",
			"attributes":        "{}",
			"droppedCount":      uint32(0),
			"flags":             plog.LogRecordFlags(0),
			"traceID":           "",
			"spanID":            "",
		},
		{
			"rSchema":           "resource_schema_url",
			"rDroppedCount":     uint32(0),
			"rAttributes":       "{}",
			"sSchema":           "scope_schema_url",
			"sName":             "",
			"sVersion":          "",
			"sDroppedCount":     uint32(0),
			"sAttributes":       "{}",
			"timestamp":         time.Unix(0, int64(0)).UTC(),
			"observedTimestamp": time.Unix(0, int64(0)).UTC(),
			"severityNumber":    "Unspecified",
			"severityText":      "",
			"body":              "log_body_2",
			"attributes":        "{}",
			"droppedCount":      uint32(0),
			"flags":             plog.LogRecordFlags(0),
			"traceID":           "",
			"spanID":            "",
		},
	}
}
