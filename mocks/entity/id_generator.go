// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	entity "github.com/mkaiho/go-lambda-api-sample/entity"
	mock "github.com/stretchr/testify/mock"
)

// IDGenerator is an autogenerated mock type for the IDGenerator type
type IDGenerator struct {
	mock.Mock
}

// From provides a mock function with given fields: value
func (_m *IDGenerator) From(value string) (entity.ID, error) {
	ret := _m.Called(value)

	var r0 entity.ID
	if rf, ok := ret.Get(0).(func(string) entity.ID); ok {
		r0 = rf(value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(entity.ID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Generate provides a mock function with given fields:
func (_m *IDGenerator) Generate() entity.ID {
	ret := _m.Called()

	var r0 entity.ID
	if rf, ok := ret.Get(0).(func() entity.ID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(entity.ID)
		}
	}

	return r0
}