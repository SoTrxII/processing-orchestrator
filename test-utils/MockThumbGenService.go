// Code generated by mockery. DO NOT EDIT.

package test_utils

import (
	mock "github.com/stretchr/testify/mock"

	thumbnail_generator "processing-orchestrator/pkg/thumb-generator/proto"
)

// MockThumbGenService is an autogenerated mock type for the ThumbGenService type
type MockThumbGenService struct {
	mock.Mock
}

type MockThumbGenService_Expecter struct {
	mock *mock.Mock
}

func (_m *MockThumbGenService) EXPECT() *MockThumbGenService_Expecter {
	return &MockThumbGenService_Expecter{mock: &_m.Mock}
}

// GenerateThumbnail provides a mock function with given fields: req
func (_m *MockThumbGenService) GenerateThumbnail(req *thumbnail_generator.ThumbnailRequest) (string, error) {
	ret := _m.Called(req)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*thumbnail_generator.ThumbnailRequest) (string, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*thumbnail_generator.ThumbnailRequest) string); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*thumbnail_generator.ThumbnailRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockThumbGenService_GenerateThumbnail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GenerateThumbnail'
type MockThumbGenService_GenerateThumbnail_Call struct {
	*mock.Call
}

// GenerateThumbnail is a helper method to define mock.On call
//   - req *thumbnail_generator.ThumbnailRequest
func (_e *MockThumbGenService_Expecter) GenerateThumbnail(req interface{}) *MockThumbGenService_GenerateThumbnail_Call {
	return &MockThumbGenService_GenerateThumbnail_Call{Call: _e.mock.On("GenerateThumbnail", req)}
}

func (_c *MockThumbGenService_GenerateThumbnail_Call) Run(run func(req *thumbnail_generator.ThumbnailRequest)) *MockThumbGenService_GenerateThumbnail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*thumbnail_generator.ThumbnailRequest))
	})
	return _c
}

func (_c *MockThumbGenService_GenerateThumbnail_Call) Return(_a0 string, _a1 error) *MockThumbGenService_GenerateThumbnail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockThumbGenService_GenerateThumbnail_Call) RunAndReturn(run func(*thumbnail_generator.ThumbnailRequest) (string, error)) *MockThumbGenService_GenerateThumbnail_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockThumbGenService creates a new instance of MockThumbGenService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockThumbGenService(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockThumbGenService {
	mock := &MockThumbGenService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
