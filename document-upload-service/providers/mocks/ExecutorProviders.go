// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	rest_service "document-upload-service/usecases/upload_service"
)

// ExecutorProviders is an autogenerated mock type for the ExecutorProviders type
type ExecutorProviders struct {
	mock.Mock
}

// GetRestServiceFactory provides a mock function with given fields:
func (_m *ExecutorProviders) GetRestServiceFactory() rest_service.RestServiceFactory {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRestServiceFactory")
	}

	var r0 rest_service.RestServiceFactory
	if rf, ok := ret.Get(0).(func() rest_service.RestServiceFactory); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rest_service.RestServiceFactory)
		}
	}

	return r0
}

// NewExecutorProviders creates a new instance of ExecutorProviders. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExecutorProviders(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExecutorProviders {
	mock := &ExecutorProviders{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
