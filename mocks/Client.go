// Code generated by mockery v2.2.1. DO NOT EDIT.

package mocks

import (
	context "context"

	sspvo "github.com/ftomza/go-sspvo"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Send provides a mock function with given fields: ctx, msg
func (_m *Client) Send(ctx context.Context, msg sspvo.Message) sspvo.Response {
	ret := _m.Called(ctx, msg)

	var r0 sspvo.Response
	if rf, ok := ret.Get(0).(func(context.Context, sspvo.Message) sspvo.Response); ok {
		r0 = rf(ctx, msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sspvo.Response)
		}
	}

	return r0
}
