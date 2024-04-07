// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	domain "github.com/Verce11o/effective-mobile-test/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// ApiCommunicator is an autogenerated mock type for the ApiCommunicator type
type ApiCommunicator struct {
	mock.Mock
}

// GetCarInfo provides a mock function with given fields: regNum
func (_m *ApiCommunicator) GetCarInfo(regNum string) (domain.Car, error) {
	ret := _m.Called(regNum)

	if len(ret) == 0 {
		panic("no return value specified for GetCarInfo")
	}

	var r0 domain.Car
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (domain.Car, error)); ok {
		return rf(regNum)
	}
	if rf, ok := ret.Get(0).(func(string) domain.Car); ok {
		r0 = rf(regNum)
	} else {
		r0 = ret.Get(0).(domain.Car)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(regNum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewApiCommunicator creates a new instance of ApiCommunicator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewApiCommunicator(t interface {
	mock.TestingT
	Cleanup(func())
}) *ApiCommunicator {
	mock := &ApiCommunicator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}