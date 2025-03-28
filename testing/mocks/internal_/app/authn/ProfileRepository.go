// Code generated by mockery v2.52.1. DO NOT EDIT.

package authn

import (
	context "context"

	web3util "github.com/growteer/api/pkg/web3util"
	mock "github.com/stretchr/testify/mock"
)

// MockProfileRepository is an autogenerated mock type for the ProfileRepository type
type MockProfileRepository struct {
	mock.Mock
}

type MockProfileRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockProfileRepository) EXPECT() *MockProfileRepository_Expecter {
	return &MockProfileRepository_Expecter{mock: &_m.Mock}
}

// Exists provides a mock function with given fields: ctx, did
func (_m *MockProfileRepository) Exists(ctx context.Context, did *web3util.DID) bool {
	ret := _m.Called(ctx, did)

	if len(ret) == 0 {
		panic("no return value specified for Exists")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *web3util.DID) bool); ok {
		r0 = rf(ctx, did)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockProfileRepository_Exists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exists'
type MockProfileRepository_Exists_Call struct {
	*mock.Call
}

// Exists is a helper method to define mock.On call
//   - ctx context.Context
//   - did *web3util.DID
func (_e *MockProfileRepository_Expecter) Exists(ctx interface{}, did interface{}) *MockProfileRepository_Exists_Call {
	return &MockProfileRepository_Exists_Call{Call: _e.mock.On("Exists", ctx, did)}
}

func (_c *MockProfileRepository_Exists_Call) Run(run func(ctx context.Context, did *web3util.DID)) *MockProfileRepository_Exists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*web3util.DID))
	})
	return _c
}

func (_c *MockProfileRepository_Exists_Call) Return(_a0 bool) *MockProfileRepository_Exists_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProfileRepository_Exists_Call) RunAndReturn(run func(context.Context, *web3util.DID) bool) *MockProfileRepository_Exists_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockProfileRepository creates a new instance of MockProfileRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockProfileRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockProfileRepository {
	mock := &MockProfileRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
