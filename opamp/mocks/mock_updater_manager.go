// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MockUpdaterManager is an autogenerated mock type for the updaterManager type
type MockUpdaterManager struct {
	mock.Mock
}

// StartAndMonitorUpdater provides a mock function with given fields:
func (_m *MockUpdaterManager) StartAndMonitorUpdater() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockUpdaterManager creates a new instance of MockUpdaterManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockUpdaterManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockUpdaterManager {
	mock := &MockUpdaterManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
