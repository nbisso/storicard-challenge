// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	context "context"

	queueclient "github.com/nbisso/storicard-challenge/infrastracture/queue_client"
	mock "github.com/stretchr/testify/mock"
)

// QueueSenderClient is an autogenerated mock type for the QueueSenderClient type
type QueueSenderClient struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *QueueSenderClient) Close() {
	_m.Called()
}

// Flush provides a mock function with given fields:
func (_m *QueueSenderClient) Flush() {
	_m.Called()
}

// InitTransaction provides a mock function with given fields: ctx
func (_m *QueueSenderClient) InitTransaction(ctx context.Context) (queueclient.QueueSenderClientTransactioner, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for InitTransaction")
	}

	var r0 queueclient.QueueSenderClientTransactioner
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (queueclient.QueueSenderClientTransactioner, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) queueclient.QueueSenderClientTransactioner); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(queueclient.QueueSenderClientTransactioner)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendMessage provides a mock function with given fields: message
func (_m *QueueSenderClient) SendMessage(message string) error {
	ret := _m.Called(message)

	if len(ret) == 0 {
		panic("no return value specified for SendMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendMessageWithHeaders provides a mock function with given fields: message, headers
func (_m *QueueSenderClient) SendMessageWithHeaders(message string, headers map[string]string) error {
	ret := _m.Called(message, headers)

	if len(ret) == 0 {
		panic("no return value specified for SendMessageWithHeaders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, map[string]string) error); ok {
		r0 = rf(message, headers)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewQueueSenderClient creates a new instance of QueueSenderClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQueueSenderClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *QueueSenderClient {
	mock := &QueueSenderClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
