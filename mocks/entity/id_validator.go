// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// IDValidator is an autogenerated mock type for the IDValidator type
type IDValidator struct {
	mock.Mock
}

// Validate provides a mock function with given fields: value
func (_m *IDValidator) Validate(value string) error {
	ret := _m.Called(value)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
