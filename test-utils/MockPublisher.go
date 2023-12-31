// Code generated by mockery. DO NOT EDIT.

package test_utils

import (
	context "context"

	client "github.com/dapr/go-sdk/client"

	mock "github.com/stretchr/testify/mock"
)

// MockPublisher is an autogenerated mock type for the Publisher type
type MockPublisher struct {
	mock.Mock
}

type MockPublisher_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPublisher) EXPECT() *MockPublisher_Expecter {
	return &MockPublisher_Expecter{mock: &_m.Mock}
}

// PublishEvent provides a mock function with given fields: ctx, pubsubName, topicName, data, opts
func (_m *MockPublisher) PublishEvent(ctx context.Context, pubsubName string, topicName string, data interface{}, opts ...client.PublishEventOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, pubsubName, topicName, data)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, interface{}, ...client.PublishEventOption) error); ok {
		r0 = rf(ctx, pubsubName, topicName, data, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockPublisher_PublishEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishEvent'
type MockPublisher_PublishEvent_Call struct {
	*mock.Call
}

// PublishEvent is a helper method to define mock.On call
//   - ctx context.Context
//   - pubsubName string
//   - topicName string
//   - data interface{}
//   - opts ...client.PublishEventOption
func (_e *MockPublisher_Expecter) PublishEvent(ctx interface{}, pubsubName interface{}, topicName interface{}, data interface{}, opts ...interface{}) *MockPublisher_PublishEvent_Call {
	return &MockPublisher_PublishEvent_Call{Call: _e.mock.On("PublishEvent",
		append([]interface{}{ctx, pubsubName, topicName, data}, opts...)...)}
}

func (_c *MockPublisher_PublishEvent_Call) Run(run func(ctx context.Context, pubsubName string, topicName string, data interface{}, opts ...client.PublishEventOption)) *MockPublisher_PublishEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]client.PublishEventOption, len(args)-4)
		for i, a := range args[4:] {
			if a != nil {
				variadicArgs[i] = a.(client.PublishEventOption)
			}
		}
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *MockPublisher_PublishEvent_Call) Return(_a0 error) *MockPublisher_PublishEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPublisher_PublishEvent_Call) RunAndReturn(run func(context.Context, string, string, interface{}, ...client.PublishEventOption) error) *MockPublisher_PublishEvent_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockPublisher creates a new instance of MockPublisher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPublisher(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPublisher {
	mock := &MockPublisher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
