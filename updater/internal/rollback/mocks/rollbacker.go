// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	action "github.com/observiq/observiq-otel-collector/updater/internal/action"
	mock "github.com/stretchr/testify/mock"
)

// Rollbacker is an autogenerated mock type for the Rollbacker type
type Rollbacker struct {
	mock.Mock
}

// AppendAction provides a mock function with given fields: _a0
func (_m *Rollbacker) AppendAction(_a0 action.RollbackableAction) {
	_m.Called(_a0)
}

// Backup provides a mock function with given fields:
func (_m *Rollbacker) Backup() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Rollback provides a mock function with given fields:
func (_m *Rollbacker) Rollback() {
	_m.Called()
}

type mockConstructorTestingTNewRollbacker interface {
	mock.TestingT
	Cleanup(func())
}

// NewRollbacker creates a new instance of Rollbacker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRollbacker(t mockConstructorTestingTNewRollbacker) *Rollbacker {
	mock := &Rollbacker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
