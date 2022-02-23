// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/mkaiho/go-lambda-api-sample/entity"
	mock "github.com/stretchr/testify/mock"
)

// UsersReader is an autogenerated mock type for the UsersReader type
type UsersReader struct {
	mock.Mock
}

// FindAll provides a mock function with given fields:
func (_m *UsersReader) FindAll() ([]entity.User, error) {
	ret := _m.Called()

	var r0 []entity.User
	if rf, ok := ret.Get(0).(func() []entity.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByID provides a mock function with given fields: id
func (_m *UsersReader) FindByID(id entity.UserID) (entity.User, error) {
	ret := _m.Called(id)

	var r0 entity.User
	if rf, ok := ret.Get(0).(func(entity.UserID) entity.User); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(entity.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(entity.UserID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
