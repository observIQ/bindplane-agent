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

package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestInitDatabaseConn(t *testing.T) {
	testCases := []struct {
		desc            string
		ctx             context.Context
		role            string
		database        string
		warehouse       string
		setExpectations func(m sqlmock.Sqlmock, role, database, warehouse string)
		expectedErr     error
	}{
		{
			desc:      "pass",
			ctx:       context.Background(),
			role:      "role",
			database:  "db",
			warehouse: "wh",
			setExpectations: func(m sqlmock.Sqlmock, role, database, warehouse string) {
				m.ExpectExec(fmt.Sprintf(`USE ROLE "%s";`, role)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS "%s";`, database)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`USE DATABASE "%s";`, database)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`USE WAREHOUSE "%s";`, warehouse)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
			},
		},
		{
			desc:      "fail USE ROLE stmt",
			ctx:       context.Background(),
			role:      "role",
			database:  "",
			warehouse: "",
			setExpectations: func(m sqlmock.Sqlmock, role, database, warehouse string) {
				m.ExpectExec(fmt.Sprintf(`USE ROLE "%s";`, role)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(fmt.Errorf("fail"))
			},
			expectedErr: fmt.Errorf("failed to call 'USE ROLE \"role\";': fail"),
		},
		{
			desc:      "fail CREATE DATABASE stmt",
			ctx:       context.Background(),
			role:      "role",
			database:  "db",
			warehouse: "",
			setExpectations: func(m sqlmock.Sqlmock, role, database, warehouse string) {
				m.ExpectExec(fmt.Sprintf(`USE ROLE "%s";`, role)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS "%s";`, database)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(fmt.Errorf("fail"))
			},
			expectedErr: fmt.Errorf("failed to create new database: fail"),
		},
		{
			desc:      "fail USE DATABASE stmt",
			ctx:       context.Background(),
			role:      "role",
			database:  "db",
			warehouse: "",
			setExpectations: func(m sqlmock.Sqlmock, role, database, warehouse string) {
				m.ExpectExec(fmt.Sprintf(`USE ROLE "%s";`, role)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS "%s";`, database)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`USE DATABASE "%s"`, database)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(fmt.Errorf("fail"))
			},
			expectedErr: fmt.Errorf("failed to call 'USE DATABASE \"db\";': fail"),
		},
		{
			desc:      "fail USE WAREHOUSE stmt",
			ctx:       context.Background(),
			role:      "role",
			database:  "db",
			warehouse: "wh",
			setExpectations: func(m sqlmock.Sqlmock, role, database, warehouse string) {
				m.ExpectExec(fmt.Sprintf(`USE ROLE "%s";`, role)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS "%s";`, database)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`USE DATABASE "%s"`, database)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`USE WAREHOUSE "%s";`, warehouse)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(fmt.Errorf("fail"))
			},
			expectedErr: fmt.Errorf("failed to call 'USE WAREHOUSE \"wh\";': fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			db, mock := NewMock(t)
			defer db.Close()

			tc.setExpectations(mock, tc.role, tc.database, tc.warehouse)

			s := &Snowflake{db: db}
			err := s.InitDatabaseConn(tc.ctx, tc.role, tc.database, tc.warehouse)

			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateSchema(t *testing.T) {
	testCases := []struct {
		desc            string
		ctx             context.Context
		schema          string
		setExpectations func(m sqlmock.Sqlmock, schema string)
		expectedErr     error
	}{
		{
			desc:   "pass",
			ctx:    context.Background(),
			schema: "schema",
			setExpectations: func(m sqlmock.Sqlmock, schema string) {
				m.ExpectExec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, schema)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`USE SCHEMA "%s";`, schema)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
			},
		},
		{
			desc:   "fail CREATE SCHEMA stmt",
			ctx:    context.Background(),
			schema: "schema",
			setExpectations: func(m sqlmock.Sqlmock, schema string) {
				m.ExpectExec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, schema)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(fmt.Errorf("fail"))
			},
			expectedErr: fmt.Errorf("failed to create schema 'schema': fail"),
		},
		{
			desc:   "fail USE SCHEMA stmt",
			ctx:    context.Background(),
			schema: "schema",
			setExpectations: func(m sqlmock.Sqlmock, schema string) {
				m.ExpectExec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s";`, schema)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectExec(fmt.Sprintf(`USE SCHEMA "%s";`, schema)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(fmt.Errorf("fail"))
			},
			expectedErr: fmt.Errorf("failed to call 'USE SCHEMA \"schema\";': fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			db, mock := NewMock(t)
			defer db.Close()

			tc.setExpectations(mock, tc.schema)

			s := &Snowflake{db: db}
			err := s.CreateSchema(tc.ctx, tc.schema)

			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateTable(t *testing.T) {
	template := `CREATE TABLE "%s";`
	testCases := []struct {
		desc            string
		ctx             context.Context
		table           string
		setExpectations func(m sqlmock.Sqlmock, table string)
		expectedErr     error
	}{
		{
			desc:  "pass",
			ctx:   context.Background(),
			table: "table",
			setExpectations: func(m sqlmock.Sqlmock, table string) {
				m.ExpectExec(fmt.Sprintf(template, table)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
			},
		},
		{
			desc: "fail CRATE TABLE stmt",
			ctx:  context.Background(),

			table: "table",
			setExpectations: func(m sqlmock.Sqlmock, table string) {
				m.ExpectExec(fmt.Sprintf(template, table)).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(fmt.Errorf("fail"))
			},
			expectedErr: fmt.Errorf("failed to create table: fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			db, mock := NewMock(t)
			defer db.Close()

			tc.setExpectations(mock, tc.table)

			s := &Snowflake{db: db}
			err := s.CreateTable(tc.ctx, tc.table, template)

			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBatchInsert(t *testing.T) {
	testCases := []struct {
		desc            string
		ctx             context.Context
		data            []map[string]any
		insert          string
		setExpectations func(m sqlmock.Sqlmock, insert string)
		expectedErr     error
	}{
		{
			desc: "pass",
			ctx:  context.Background(),
			data: []map[string]any{
				{
					"map1key1": "1",
					"map1key2": "2",
				},
			},
			insert: "insert",
			setExpectations: func(m sqlmock.Sqlmock, insert string) {
				m.ExpectBegin().WillReturnError(nil)
				m.ExpectExec(insert).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
				m.ExpectCommit()
			},
		},
		{
			desc:   "fail BeginTxx",
			ctx:    context.Background(),
			data:   []map[string]any{},
			insert: "insert",
			setExpectations: func(m sqlmock.Sqlmock, insert string) {
				m.ExpectBegin().WillReturnError(fmt.Errorf("fail"))
			},
			expectedErr: fmt.Errorf("failed to create transaction: fail"),
		},
		{
			desc: "fail insert stmt",
			ctx:  context.Background(),
			data: []map[string]any{
				{
					"map1key1": "1",
					"map1key2": "2",
				},
			},
			insert: "insert",
			setExpectations: func(m sqlmock.Sqlmock, insert string) {
				m.ExpectBegin().WillReturnError(nil)
				m.ExpectExec(insert).WithoutArgs().WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(fmt.Errorf("fail"))
				m.ExpectRollback()
			},
			expectedErr: fmt.Errorf("failed to execute batch insert: fail"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			db, mock := NewMock(t)
			defer db.Close()

			tc.setExpectations(mock, tc.insert)

			s := &Snowflake{db: db}
			err := s.BatchInsert(tc.ctx, tc.data, tc.insert)

			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestClose(t *testing.T) {
	// nil db case
	s := &Snowflake{}
	require.NoError(t, s.Close())

	// mock returns err
	db, mock := NewMock(t)
	mock.ExpectClose().WillReturnError(fmt.Errorf("fail"))
	s.db = db
	require.ErrorContains(t, s.Close(), "fail")
}

func NewMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock
}
