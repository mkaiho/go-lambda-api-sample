// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// UserID is an autogenerated mock type for the UserID type
type UserID struct {
	mock.Mock
}

// Value provides a mock function with given fields:
func (_m *UserID) Value() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}