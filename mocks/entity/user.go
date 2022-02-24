// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/mkaiho/go-lambda-api-sample/entity"
	mock "github.com/stretchr/testify/mock"
)

// User is an autogenerated mock type for the User type
type User struct {
	mock.Mock
}

// Name provides a mock function with given fields:
func (_m *User) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// UserID provides a mock function with given fields:
func (_m *User) UserID() entity.UserID {
	ret := _m.Called()

	var r0 entity.UserID
	if rf, ok := ret.Get(0).(func() entity.UserID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(entity.UserID)
		}
	}

	return r0
}