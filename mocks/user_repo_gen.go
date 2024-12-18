// Code generated by mockery v2.49.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/moonicy/goph-keeper-yandex/internal/entity"
	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

type UserRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepository) EXPECT() *UserRepository_Expecter {
	return &UserRepository_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, user
func (_m *UserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) (entity.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, entity.User) entity.User); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, entity.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type UserRepository_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - user entity.User
func (_e *UserRepository_Expecter) Create(ctx interface{}, user interface{}) *UserRepository_Create_Call {
	return &UserRepository_Create_Call{Call: _e.mock.On("Create", ctx, user)}
}

func (_c *UserRepository_Create_Call) Run(run func(ctx context.Context, user entity.User)) *UserRepository_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.User))
	})
	return _c
}

func (_c *UserRepository_Create_Call) Return(_a0 entity.User, _a1 error) *UserRepository_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Create_Call) RunAndReturn(run func(context.Context, entity.User) (entity.User, error)) *UserRepository_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, login
func (_m *UserRepository) Get(ctx context.Context, login string) (entity.User, error) {
	ret := _m.Called(ctx, login)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 entity.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (entity.User, error)); ok {
		return rf(ctx, login)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.User); ok {
		r0 = rf(ctx, login)
	} else {
		r0 = ret.Get(0).(entity.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepository_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type UserRepository_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - login string
func (_e *UserRepository_Expecter) Get(ctx interface{}, login interface{}) *UserRepository_Get_Call {
	return &UserRepository_Get_Call{Call: _e.mock.On("Get", ctx, login)}
}

func (_c *UserRepository_Get_Call) Run(run func(ctx context.Context, login string)) *UserRepository_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepository_Get_Call) Return(_a0 entity.User, _a1 error) *UserRepository_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepository_Get_Call) RunAndReturn(run func(context.Context, string) (entity.User, error)) *UserRepository_Get_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
