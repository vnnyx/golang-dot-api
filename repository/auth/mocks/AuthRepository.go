// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "github.com/vnnyx/golang-dot-api/model"
)

// AuthRepository is an autogenerated mock type for the AuthRepository type
type AuthRepository struct {
	mock.Mock
}

// DeleteToken provides a mock function with given fields: ctx, accessUuid
func (_m *AuthRepository) DeleteToken(ctx context.Context, accessUuid string) error {
	ret := _m.Called(ctx, accessUuid)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, accessUuid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FlushAll provides a mock function with given fields: ctx
func (_m *AuthRepository) FlushAll(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetToken provides a mock function with given fields: ctx, accessUuid
func (_m *AuthRepository) GetToken(ctx context.Context, accessUuid string) (string, error) {
	ret := _m.Called(ctx, accessUuid)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, accessUuid)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, accessUuid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreToken provides a mock function with given fields: ctx, details
func (_m *AuthRepository) StoreToken(ctx context.Context, details model.TokenDetails) error {
	ret := _m.Called(ctx, details)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.TokenDetails) error); ok {
		r0 = rf(ctx, details)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewAuthRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewAuthRepository creates a new instance of AuthRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAuthRepository(t mockConstructorTestingTNewAuthRepository) *AuthRepository {
	mock := &AuthRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
