// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	web "github.com/mkaiho/go-lambda-api-sample/adapter/web"
	mock "github.com/stretchr/testify/mock"
)

// GetUserHandler is an autogenerated mock type for the GetUserHandler type
type GetUserHandler struct {
	mock.Mock
}

// Handle provides a mock function with given fields: req
func (_m *GetUserHandler) Handle(req web.GetUserRequest) *web.GetUserResponse {
	ret := _m.Called(req)

	var r0 *web.GetUserResponse
	if rf, ok := ret.Get(0).(func(web.GetUserRequest) *web.GetUserResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*web.GetUserResponse)
		}
	}

	return r0
}
