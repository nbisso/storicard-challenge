// Code generated by mockery v2.45.0. DO NOT EDIT.

package mocks

import (
	kafka "github.com/confluentinc/confluent-kafka-go/kafka"
	mock "github.com/stretchr/testify/mock"
)

// QueueConsumerClient is an autogenerated mock type for the QueueConsumerClient type
type QueueConsumerClient struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *QueueConsumerClient) Close() {
	_m.Called()
}

// CommitMessage provides a mock function with given fields: m
func (_m *QueueConsumerClient) CommitMessage(m *kafka.Message) error {
	ret := _m.Called(m)

	if len(ret) == 0 {
		panic("no return value specified for CommitMessage")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*kafka.Message) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReadMessage provides a mock function with given fields:
func (_m *QueueConsumerClient) ReadMessage() (*kafka.Message, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReadMessage")
	}

	var r0 *kafka.Message
	var r1 error
	if rf, ok := ret.Get(0).(func() (*kafka.Message, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *kafka.Message); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*kafka.Message)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewQueueConsumerClient creates a new instance of QueueConsumerClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQueueConsumerClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *QueueConsumerClient {
	mock := &QueueConsumerClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
