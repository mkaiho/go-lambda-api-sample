// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	web "github.com/mkaiho/go-lambda-api-sample/adapter/web"
	mock "github.com/stretchr/testify/mock"
)

// DeleteUserHandler is an autogenerated mock type for the DeleteUserHandler type
type DeleteUserHandler struct {
	mock.Mock
}

// Handle provides a mock function with given fields: req
func (_m *DeleteUserHandler) Handle(req web.DeleteUserRequest) *web.DeleteUserResponse {
	ret := _m.Called(req)

	var r0 *web.DeleteUserResponse
	if rf, ok := ret.Get(0).(func(web.DeleteUserRequest) *web.DeleteUserResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*web.DeleteUserResponse)
		}
	}

	return r0
}
