// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/vnnyx/golang-dot-api/model/entity"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// TransactionRepository is an autogenerated mock type for the TransactionRepository type
type TransactionRepository struct {
	mock.Mock
}

// DeleteAllTransaction provides a mock function with given fields: ctx
func (_m *TransactionRepository) DeleteAllTransaction(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTransaction provides a mock function with given fields: ctx, transactionId
func (_m *TransactionRepository) DeleteTransaction(ctx context.Context, transactionId string) error {
	ret := _m.Called(ctx, transactionId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, transactionId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTransactionByUserId provides a mock function with given fields: ctx, tx, userId
func (_m *TransactionRepository) DeleteTransactionByUserId(ctx context.Context, tx *gorm.DB, userId string) error {
	ret := _m.Called(ctx, tx, userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string) error); ok {
		r0 = rf(ctx, tx, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAllTransaction provides a mock function with given fields: ctx
func (_m *TransactionRepository) FindAllTransaction(ctx context.Context) ([]entity.Transaction, error) {
	ret := _m.Called(ctx)

	var r0 []entity.Transaction
	if rf, ok := ret.Get(0).(func(context.Context) []entity.Transaction); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTransactionByID provides a mock function with given fields: ctx, transactionId
func (_m *TransactionRepository) FindTransactionByID(ctx context.Context, transactionId string) (entity.Transaction, error) {
	ret := _m.Called(ctx, transactionId)

	var r0 entity.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string) entity.Transaction); ok {
		r0 = rf(ctx, transactionId)
	} else {
		r0 = ret.Get(0).(entity.Transaction)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, transactionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTransactionByUserId provides a mock function with given fields: ctx, userId
func (_m *TransactionRepository) FindTransactionByUserId(ctx context.Context, userId string) ([]entity.Transaction, error) {
	ret := _m.Called(ctx, userId)

	var r0 []entity.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string) []entity.Transaction); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertTransaction provides a mock function with given fields: ctx, _a1
func (_m *TransactionRepository) InsertTransaction(ctx context.Context, _a1 entity.Transaction) (entity.Transaction, error) {
	ret := _m.Called(ctx, _a1)

	var r0 entity.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, entity.Transaction) entity.Transaction); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(entity.Transaction)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, entity.Transaction) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTransaction provides a mock function with given fields: ctx, _a1
func (_m *TransactionRepository) UpdateTransaction(ctx context.Context, _a1 entity.Transaction) (entity.Transaction, error) {
	ret := _m.Called(ctx, _a1)

	var r0 entity.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, entity.Transaction) entity.Transaction); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(entity.Transaction)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, entity.Transaction) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTransactionRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionRepository creates a new instance of TransactionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionRepository(t mockConstructorTestingTNewTransactionRepository) *TransactionRepository {
	mock := &TransactionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
