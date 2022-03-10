// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/mkaiho/go-lambda-api-sample/entity"
	mock "github.com/stretchr/testify/mock"
)

// UsersWriter is an autogenerated mock type for the UsersWriter type
type UsersWriter struct {
	mock.Mock
}

// Delete provides a mock function with given fields: id
func (_m *UsersWriter) Delete(id entity.UserID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(entity.UserID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Insert provides a mock function with given fields: user
func (_m *UsersWriter) Insert(user entity.User) (entity.User, error) {
	ret := _m.Called(user)

	var r0 entity.User
	if rf, ok := ret.Get(0).(func(entity.User) entity.User); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(entity.User) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
