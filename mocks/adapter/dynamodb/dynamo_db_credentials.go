// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DynamoDBCredentials is an autogenerated mock type for the DynamoDBCredentials type
type DynamoDBCredentials struct {
	mock.Mock
}

// AccessKeyID provides a mock function with given fields:
func (_m *DynamoDBCredentials) AccessKeyID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SecretAccessKey provides a mock function with given fields:
func (_m *DynamoDBCredentials) SecretAccessKey() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}