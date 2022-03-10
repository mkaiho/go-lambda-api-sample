// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	mock "github.com/stretchr/testify/mock"
)

// DynamoDBClient is an autogenerated mock type for the DynamoDBClient type
type DynamoDBClient struct {
	mock.Mock
}

// DeleteItem provides a mock function with given fields: input
func (_m *DynamoDBClient) DeleteItem(input dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	ret := _m.Called(input)

	var r0 *dynamodb.DeleteItemOutput
	if rf, ok := ret.Get(0).(func(dynamodb.DeleteItemInput) *dynamodb.DeleteItemOutput); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.DeleteItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dynamodb.DeleteItemInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetItem provides a mock function with given fields: input
func (_m *DynamoDBClient) GetItem(input dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	ret := _m.Called(input)

	var r0 *dynamodb.GetItemOutput
	if rf, ok := ret.Get(0).(func(dynamodb.GetItemInput) *dynamodb.GetItemOutput); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.GetItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dynamodb.GetItemInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutItem provides a mock function with given fields: input
func (_m *DynamoDBClient) PutItem(input dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	ret := _m.Called(input)

	var r0 *dynamodb.PutItemOutput
	if rf, ok := ret.Get(0).(func(dynamodb.PutItemInput) *dynamodb.PutItemOutput); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.PutItemOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dynamodb.PutItemInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: input
func (_m *DynamoDBClient) Query(input dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	ret := _m.Called(input)

	var r0 *dynamodb.QueryOutput
	if rf, ok := ret.Get(0).(func(dynamodb.QueryInput) *dynamodb.QueryOutput); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.QueryOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dynamodb.QueryInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Scan provides a mock function with given fields: input
func (_m *DynamoDBClient) Scan(input dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	ret := _m.Called(input)

	var r0 *dynamodb.ScanOutput
	if rf, ok := ret.Get(0).(func(dynamodb.ScanInput) *dynamodb.ScanOutput); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dynamodb.ScanOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(dynamodb.ScanInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
