// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockDatabase is an autogenerated mock type for the Database type
type MockDatabase struct {
	mock.Mock
}

// BatchInsert provides a mock function with given fields: ctx, data, sql
func (_m *MockDatabase) BatchInsert(ctx context.Context, data []map[string]interface{}, sql string) error {
	ret := _m.Called(ctx, data, sql)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []map[string]interface{}, string) error); ok {
		r0 = rf(ctx, data, sql)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *MockDatabase) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateSchema provides a mock function with given fields: ctx, schema
func (_m *MockDatabase) CreateSchema(ctx context.Context, schema string) error {
	ret := _m.Called(ctx, schema)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, schema)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateTable provides a mock function with given fields: ctx, sql
func (_m *MockDatabase) CreateTable(ctx context.Context, sql string) error {
	ret := _m.Called(ctx, sql)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, sql)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InitDatabaseConn provides a mock function with given fields: ctx, roles
func (_m *MockDatabase) InitDatabaseConn(ctx context.Context, roles string) error {
	ret := _m.Called(ctx, roles)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, roles)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockDatabase creates a new instance of MockDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabase {
	mock := &MockDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}