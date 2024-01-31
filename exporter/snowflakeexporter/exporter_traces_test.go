package snowflakeexporter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/database"
	"github.com/observiq/bindplane-agent/exporter/snowflakeexporter/internal/database/mocks"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func TestNewTracesExporter(t *testing.T) {
	c := &Config{
		AccountIdentifier: "id",
		Username:          "user",
		Password:          "pass",
		Database:          "db",
		Traces: &TelemetryConfig{
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
			expectedErr: fmt.Errorf("failed to create new database connection for traces: fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			exp, err := newTracesExporter(
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

func TestTracesCapabilities(t *testing.T) {
	e := &tracesExporter{}
	c := e.Capabilities()
	require.False(t, c.MutatesData)
}

func TestTracesStart(t *testing.T) {
	c := &Config{
		AccountIdentifier: "id",
		Username:          "user",
		Password:          "pass",
		Database:          "db",
		Traces: &TelemetryConfig{
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
				m.On("CreateSchema", ctx, c.Traces.Schema).Return(nil)
				m.On("CreateTable", ctx, c.Database, c.Traces.Schema, c.Traces.Table, createTracesTableSnowflakeTemplate).Return(nil)
				return m
			},
		},
		{
			desc: "Fail CreateSchema",
			ctx:  context.Background(),
			mockGen: func(t *testing.T, ctx context.Context, c *Config) *mocks.MockDatabase {
				m := mocks.NewMockDatabase(t)
				m.On("CreateSchema", ctx, c.Traces.Schema).Return(fmt.Errorf("fail"))
				return m
			},
			expectedErr: fmt.Errorf("failed to create traces schema: fail"),
		},
		{
			desc: "Fail CreateTable",
			ctx:  context.Background(),
			mockGen: func(t *testing.T, ctx context.Context, c *Config) *mocks.MockDatabase {
				m := mocks.NewMockDatabase(t)
				m.On("CreateSchema", ctx, c.Traces.Schema).Return(nil)
				m.On("CreateTable", ctx, c.Database, c.Traces.Schema, c.Traces.Table, createTracesTableSnowflakeTemplate).Return(fmt.Errorf("fail"))
				return m
			},
			expectedErr: fmt.Errorf("failed to create traces table: fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tracesExp, err := newTracesExporter(
				tc.ctx,
				c,
				exportertest.NewNopCreateSettings(),
				func(ctx context.Context, dsn string) (database.Database, error) { return nil, nil },
			)
			require.NoError(t, err)
			tracesExp.db = tc.mockGen(t, tc.ctx, c)

			err = tracesExp.start(tc.ctx, nil)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}

func TestTracesShutdown(t *testing.T) {
	// no db
	ctx := context.Background()
	e := &tracesExporter{}
	require.NoError(t, e.shutdown(ctx))

	// db & error
	mock := mocks.NewMockDatabase(t)
	mock.On("Close").Return(fmt.Errorf("fail")).Once()
	e.db = mock

	require.ErrorContains(t, e.shutdown(ctx), "fail")
}

func TestTracesDataPusher(t *testing.T) {
	c := &Config{
		AccountIdentifier: "id",
		Username:          "user",
		Password:          "pass",
		Database:          "db",
		Traces: &TelemetryConfig{
			Schema: "schema",
			Table:  "table",
		},
	}

	testCases := []struct {
		desc        string
		ctx         context.Context
		traceGen    func(t *testing.T) ptrace.Traces
		mapGen      func(t *testing.T) []map[string]any
		mockGen     func(t *testing.T, ctx context.Context, warehouse, sql string, data []map[string]any) *mocks.MockDatabase
		expectedErr error
	}{
		{
			desc:     "Simple pass",
			ctx:      context.Background(),
			traceGen: generateTraceData1,
			mapGen:   generateTraceMaps1,
			mockGen: func(t *testing.T, ctx context.Context, warehouse, sql string, data []map[string]any) *mocks.MockDatabase {
				m := mocks.NewMockDatabase(t)
				m.On("BatchInsert", ctx, data, warehouse, sql).Return(nil)
				return m
			},
		},
		{
			desc:     "Simple fail",
			ctx:      context.Background(),
			traceGen: generateTraceData1,
			mapGen:   generateTraceMaps1,
			mockGen: func(t *testing.T, ctx context.Context, warehouse, sql string, data []map[string]any) *mocks.MockDatabase {
				m := mocks.NewMockDatabase(t)
				m.On("BatchInsert", ctx, data, warehouse, sql).Return(fmt.Errorf("fail"))
				return m
			},
			expectedErr: fmt.Errorf("failed to insert trace data: fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tracesExp, err := newTracesExporter(
				tc.ctx,
				c,
				exportertest.NewNopCreateSettings(),
				func(ctx context.Context, dsn string) (database.Database, error) { return nil, nil },
			)
			require.NoError(t, err)
			tracesExp.db = tc.mockGen(t, tc.ctx, c.Warehouse, tracesExp.insertSQL, tc.mapGen(t))

			err = tracesExp.tracesDataPusher(tc.ctx, tc.traceGen(t))
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}

func generateTraceData1(t *testing.T) ptrace.Traces {
	traces := ptrace.NewTraces()
	rSpan := traces.ResourceSpans().AppendEmpty()
	rSpan.SetSchemaUrl("resource_schema_url")
	sSpan := rSpan.ScopeSpans().AppendEmpty()
	sSpan.SetSchemaUrl("scope_schema_url")
	for i := 0; i < 3; i++ {
		s := sSpan.Spans().AppendEmpty()
		s.SetName(fmt.Sprintf("span_%d", i))
	}
	return traces
}

func generateTraceMaps1(t *testing.T) []map[string]any {
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
			"traceID":           "",
			"spanID":            "",
			"traceState":        "",
			"parentSpanID":      "",
			"name":              "span_0",
			"kind":              "Unspecified",
			"startTime":         time.Unix(0, int64(0)).UTC(),
			"endTime":           time.Unix(0, int64(0)).UTC(),
			"droppedCount":      uint32(0),
			"attributes":        "{}",
			"statusMessage":     "",
			"statusCode":        "Unset",
			"eventTimes":        pq.StringArray{},
			"eventNames":        pq.StringArray{},
			"eventDroppedCount": pq.Int32Array{},
			"eventAttributes":   pq.StringArray{},
			"linkTraceIDs":      pq.StringArray{},
			"linkSpanIDs":       pq.StringArray{},
			"linkTraceStates":   pq.StringArray{},
			"linkDroppedCount":  pq.Int32Array{},
			"linkAttributes":    pq.StringArray{},
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
			"traceID":           "",
			"spanID":            "",
			"traceState":        "",
			"parentSpanID":      "",
			"name":              "span_1",
			"kind":              "Unspecified",
			"startTime":         time.Unix(0, int64(0)).UTC(),
			"endTime":           time.Unix(0, int64(0)).UTC(),
			"droppedCount":      uint32(0),
			"attributes":        "{}",
			"statusMessage":     "",
			"statusCode":        "Unset",
			"eventTimes":        pq.StringArray{},
			"eventNames":        pq.StringArray{},
			"eventDroppedCount": pq.Int32Array{},
			"eventAttributes":   pq.StringArray{},
			"linkTraceIDs":      pq.StringArray{},
			"linkSpanIDs":       pq.StringArray{},
			"linkTraceStates":   pq.StringArray{},
			"linkDroppedCount":  pq.Int32Array{},
			"linkAttributes":    pq.StringArray{},
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
			"traceID":           "",
			"spanID":            "",
			"traceState":        "",
			"parentSpanID":      "",
			"name":              "span_2",
			"kind":              "Unspecified",
			"startTime":         time.Unix(0, int64(0)).UTC(),
			"endTime":           time.Unix(0, int64(0)).UTC(),
			"droppedCount":      uint32(0),
			"attributes":        "{}",
			"statusMessage":     "",
			"statusCode":        "Unset",
			"eventTimes":        pq.StringArray{},
			"eventNames":        pq.StringArray{},
			"eventDroppedCount": pq.Int32Array{},
			"eventAttributes":   pq.StringArray{},
			"linkTraceIDs":      pq.StringArray{},
			"linkSpanIDs":       pq.StringArray{},
			"linkTraceStates":   pq.StringArray{},
			"linkDroppedCount":  pq.Int32Array{},
			"linkAttributes":    pq.StringArray{},
		},
	}
}
