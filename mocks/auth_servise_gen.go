// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/moonicy/goph-keeper-yandex/internal/entity"

	mock "github.com/stretchr/testify/mock"
)

// AuthService is an autogenerated mock type for the AuthService type
type AuthService struct {
	mock.Mock
}

type AuthService_Expecter struct {
	mock *mock.Mock
}

func (_m *AuthService) EXPECT() *AuthService_Expecter {
	return &AuthService_Expecter{mock: &_m.Mock}
}

// Login provides a mock function with given fields: ctx, login, password
func (_m *AuthService) Login(ctx context.Context, login string, password string) (string, string, error) {
	ret := _m.Called(ctx, login, password)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, string, error)); ok {
		return rf(ctx, login, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, login, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) string); ok {
		r1 = rf(ctx, login, password)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, string, string) error); ok {
		r2 = rf(ctx, login, password)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// AuthService_Login_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Login'
type AuthService_Login_Call struct {
	*mock.Call
}

// Login is a helper method to define mock.On call
//   - ctx context.Context
//   - login string
//   - password string
func (_e *AuthService_Expecter) Login(ctx interface{}, login interface{}, password interface{}) *AuthService_Login_Call {
	return &AuthService_Login_Call{Call: _e.mock.On("Login", ctx, login, password)}
}

func (_c *AuthService_Login_Call) Run(run func(ctx context.Context, login string, password string)) *AuthService_Login_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *AuthService_Login_Call) Return(_a0 string, _a1 string, _a2 error) *AuthService_Login_Call {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *AuthService_Login_Call) RunAndReturn(run func(context.Context, string, string) (string, string, error)) *AuthService_Login_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function with given fields: ctx, login, password
func (_m *AuthService) Register(ctx context.Context, login string, password string) (entity.User, error) {
	ret := _m.Called(ctx, login, password)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (entity.User, error)); ok {
		return rf(ctx, login, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) entity.User); ok {
		r0 = rf(ctx, login, password)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, login, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AuthService_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type AuthService_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - ctx context.Context
//   - login string
//   - password string
func (_e *AuthService_Expecter) Register(ctx interface{}, login interface{}, password interface{}) *AuthService_Register_Call {
	return &AuthService_Register_Call{Call: _e.mock.On("Register", ctx, login, password)}
}

func (_c *AuthService_Register_Call) Run(run func(ctx context.Context, login string, password string)) *AuthService_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *AuthService_Register_Call) Return(_a0 entity.User, _a1 error) *AuthService_Register_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AuthService_Register_Call) RunAndReturn(run func(context.Context, string, string) (entity.User, error)) *AuthService_Register_Call {
	_c.Call.Return(run)
	return _c
}

// NewAuthService creates a new instance of AuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthService(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthService {
	mock := &AuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
