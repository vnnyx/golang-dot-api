package unit

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	mockTransactionRepository "github.com/vnnyx/golang-dot-api/repository/transaction/mocks"
	mockUserRepository "github.com/vnnyx/golang-dot-api/repository/user/mocks"
	"github.com/vnnyx/golang-dot-api/service/transaction"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		req web.TransactionCreateRequest
	}
	type mockFindUserByIDRepository struct {
		res entity.User
		err error
	}
	type mockInsertTransactionRepository struct {
		res entity.Transaction
		err error
	}
	tests := []struct {
		name                            string
		args                            args
		mockFindUserByIDRepository      *mockFindUserByIDRepository
		mockInsertTransactionRepository *mockInsertTransactionRepository
		want                            web.TransactionResponse
		wantErr                         bool
	}{
		{
			name: "Transaction CreateTransaction Success",
			args: args{
				ctx: context.TODO(),
				req: web.TransactionCreateRequest{
					Name:   "product_test",
					UserID: "123",
				},
			},
			mockFindUserByIDRepository: &mockFindUserByIDRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
					Password:  "password",
				},
				err: nil,
			},
			mockInsertTransactionRepository: &mockInsertTransactionRepository{
				res: entity.Transaction{
					TransactionID: "456",
					Name:          "product_test",
					UserID:        "123",
				},
				err: nil,
			},
			want: web.TransactionResponse{
				TransactionID: "456",
				Name:          "product_test",
				UserID:        "123",
			},
			wantErr: false,
		},
		{
			name: "Error When Find User By ID",
			args: args{
				ctx: context.TODO(),
				req: web.TransactionCreateRequest{
					Name:   "product_test",
					UserID: "1234",
				},
			},
			mockFindUserByIDRepository: &mockFindUserByIDRepository{
				res: entity.User{},
				err: errors.New("user not found"),
			},
			want:    web.TransactionResponse{},
			wantErr: true,
		},
		{
			name: "Error When Insert data to DB",
			args: args{
				ctx: context.TODO(),
				req: web.TransactionCreateRequest{
					Name:   "product_test",
					UserID: "123",
				},
			},
			mockFindUserByIDRepository: &mockFindUserByIDRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
					Password:  "password",
				},
				err: nil,
			},
			mockInsertTransactionRepository: &mockInsertTransactionRepository{
				res: entity.Transaction{},
				err: errors.New("error"),
			},
			want:    web.TransactionResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)

			if tt.mockFindUserByIDRepository != nil {
				mockUserRepository.On("FindUserByID", tt.args.ctx, mock.Anything).Return(tt.mockFindUserByIDRepository.res, tt.mockFindUserByIDRepository.err)
			}
			if tt.mockInsertTransactionRepository != nil {
				mockTransactionRepository.On("InsertTransaction", tt.args.ctx, mock.Anything).Return(tt.mockInsertTransactionRepository.res, tt.mockInsertTransactionRepository.err)
			}

			transactionId := gomonkey.ApplyFunc(uuid.NewString, func() string {
				return "456"
			})
			defer transactionId.Reset()

			transactionService := transaction.NewTransactionService(mockTransactionRepository, mockUserRepository)
			got, err := transactionService.CreateTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.CreateTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.CreateTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionService_GetTransactionById(t *testing.T) {
	type args struct {
		ctx context.Context
		req string
	}
	type mockFindTransactionByIDRepository struct {
		res entity.Transaction
		err error
	}
	tests := []struct {
		name                              string
		args                              args
		mockFindTransactionByIDRepository *mockFindTransactionByIDRepository
		want                              web.TransactionResponse
		wantErr                           bool
	}{
		{
			name: "Transaction CreateTransaction Success",
			args: args{
				ctx: context.TODO(),
				req: "456",
			},
			mockFindTransactionByIDRepository: &mockFindTransactionByIDRepository{
				res: entity.Transaction{
					TransactionID: "456",
					Name:          "product_test",
					UserID:        "123",
				},
				err: nil,
			},
			want: web.TransactionResponse{
				TransactionID: "456",
				Name:          "product_test",
				UserID:        "123",
			},
			wantErr: false,
		},
		{
			name: "Transaction Not Found",
			args: args{
				ctx: context.TODO(),
				req: "4567",
			},
			mockFindTransactionByIDRepository: &mockFindTransactionByIDRepository{
				res: entity.Transaction{},
				err: errors.New("transaction not found"),
			},
			want:    web.TransactionResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)

			if tt.mockFindTransactionByIDRepository != nil {
				mockTransactionRepository.On("FindTransactionByID", tt.args.ctx, mock.Anything).Return(tt.mockFindTransactionByIDRepository.res, tt.mockFindTransactionByIDRepository.err)
			}

			transactionId := gomonkey.ApplyFunc(uuid.NewString, func() string {
				return "456"
			})
			defer transactionId.Reset()

			transactionService := transaction.NewTransactionService(mockTransactionRepository, mockUserRepository)
			got, err := transactionService.GetTransactionById(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetTransactionById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetTransactionById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionService_GetAllTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type mockFindAllTransactionRepository struct {
		res []entity.Transaction
		err error
	}
	tests := []struct {
		name                             string
		args                             args
		mockFindAllTransactionRepository *mockFindAllTransactionRepository
		want                             []web.TransactionResponse
		wantErr                          bool
	}{
		{
			name: "Transaction CreateTransaction Success",
			args: args{
				ctx: context.TODO(),
			},
			mockFindAllTransactionRepository: &mockFindAllTransactionRepository{
				res: []entity.Transaction{
					{
						TransactionID: "456",
						Name:          "product_test",
						UserID:        "123",
					},
				},
				err: nil,
			},
			want: []web.TransactionResponse{
				{
					TransactionID: "456",
					Name:          "product_test",
					UserID:        "123",
				},
			},
			wantErr: false,
		},
		{
			name: "Error When Getting Data From DB",
			args: args{
				ctx: context.TODO(),
			},
			mockFindAllTransactionRepository: &mockFindAllTransactionRepository{
				res: []entity.Transaction{},
				err: errors.New("error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)

			if tt.mockFindAllTransactionRepository != nil {
				mockTransactionRepository.On("FindAllTransaction", tt.args.ctx, mock.Anything).Return(tt.mockFindAllTransactionRepository.res, tt.mockFindAllTransactionRepository.err)
			}

			transactionId := gomonkey.ApplyFunc(uuid.NewString, func() string {
				return "456"
			})
			defer transactionId.Reset()

			transactionService := transaction.NewTransactionService(mockTransactionRepository, mockUserRepository)
			got, err := transactionService.GetAllTransaction(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetAllTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetAllTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionService_GetTransactionByUserId(t *testing.T) {
	type args struct {
		ctx context.Context
		req string
	}
	type mockFindUserByIDRepository struct {
		res entity.User
		err error
	}
	type mockFindTransactionByUserIdRepository struct {
		res []entity.Transaction
		err error
	}
	tests := []struct {
		name                                  string
		args                                  args
		mockFindUserByIDRepository            *mockFindUserByIDRepository
		mockFindTransactionByUserIdRepository *mockFindTransactionByUserIdRepository
		want                                  []web.TransactionResponse
		wantErr                               bool
	}{
		{
			name: "GetTransaction By User ID Success",
			args: args{
				ctx: context.TODO(),
				req: "123",
			},
			mockFindUserByIDRepository: &mockFindUserByIDRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
					Password:  "password",
				},
				err: nil,
			},
			mockFindTransactionByUserIdRepository: &mockFindTransactionByUserIdRepository{
				res: []entity.Transaction{
					{
						TransactionID: "456",
						Name:          "product_test",
						UserID:        "123",
					},
				},
				err: nil,
			},
			want: []web.TransactionResponse{
				{
					TransactionID: "456",
					Name:          "product_test",
					UserID:        "123",
				},
			},
			wantErr: false,
		},
		{
			name: "Error When Find User ID",
			args: args{
				ctx: context.TODO(),
				req: "1234",
			},
			mockFindUserByIDRepository: &mockFindUserByIDRepository{
				res: entity.User{},
				err: errors.New("user not found"),
			},
			wantErr: true,
		},
		{
			name: "Error When Get Transaction Data",
			args: args{
				ctx: context.TODO(),
				req: "123",
			},
			mockFindUserByIDRepository: &mockFindUserByIDRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
					Password:  "password",
				},
				err: nil,
			},
			mockFindTransactionByUserIdRepository: &mockFindTransactionByUserIdRepository{
				res: []entity.Transaction{},
				err: errors.New("error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)

			if tt.mockFindUserByIDRepository != nil {
				mockUserRepository.On("FindUserByID", tt.args.ctx, mock.Anything).Return(tt.mockFindUserByIDRepository.res, tt.mockFindUserByIDRepository.err)
			}
			if tt.mockFindTransactionByUserIdRepository != nil {
				mockTransactionRepository.On("FindTransactionByUserId", tt.args.ctx, mock.Anything).Return(tt.mockFindTransactionByUserIdRepository.res, tt.mockFindTransactionByUserIdRepository.err)
			}

			transactionService := transaction.NewTransactionService(mockTransactionRepository, mockUserRepository)
			got, err := transactionService.GetTransactionByUserId(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetTransactionByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(reflect.TypeOf(got))
			fmt.Println(reflect.TypeOf(tt.want))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetTransactionByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionService_UpdateTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		req web.TransactionUpdateRequest
	}
	type mockFindTransactionByIDRepository struct {
		res entity.Transaction
		err error
	}
	type mockUpdateTransactionRepository struct {
		res entity.Transaction
		err error
	}
	tests := []struct {
		name                              string
		args                              args
		mockFindTransactionByIDRepository *mockFindTransactionByIDRepository
		mockUpdateTransactionRepository   *mockUpdateTransactionRepository
		want                              web.TransactionResponse
		wantErr                           bool
	}{
		{
			name: "Update Transaction Success",
			args: args{
				ctx: context.TODO(),
				req: web.TransactionUpdateRequest{
					TransactionID: "456",
					Name:          "product_test_update",
				},
			},
			mockFindTransactionByIDRepository: &mockFindTransactionByIDRepository{
				res: entity.Transaction{
					TransactionID: "456",
					Name:          "product_test",
					UserID:        "123",
				},
				err: nil,
			},
			mockUpdateTransactionRepository: &mockUpdateTransactionRepository{
				res: entity.Transaction{
					TransactionID: "456",
					Name:          "product_test_update",
					UserID:        "123",
				},
				err: nil,
			},
			want: web.TransactionResponse{
				TransactionID: "456",
				Name:          "product_test_update",
				UserID:        "123",
			},
			wantErr: false,
		},
		{
			name: "Transaction Not Found",
			args: args{
				ctx: context.TODO(),
				req: web.TransactionUpdateRequest{
					TransactionID: "4567",
					Name:          "product_test_update",
				},
			},
			mockFindTransactionByIDRepository: &mockFindTransactionByIDRepository{
				res: entity.Transaction{},
				err: errors.New("transaction not found"),
			},
			want:    web.TransactionResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)

			if tt.mockFindTransactionByIDRepository != nil {
				mockTransactionRepository.On("FindTransactionByID", tt.args.ctx, mock.Anything).Return(tt.mockFindTransactionByIDRepository.res, tt.mockFindTransactionByIDRepository.err)
			}
			if tt.mockUpdateTransactionRepository != nil {
				mockTransactionRepository.On("UpdateTransaction", tt.args.ctx, mock.Anything).Return(tt.mockUpdateTransactionRepository.res, tt.mockUpdateTransactionRepository.err)
			}

			transactionId := gomonkey.ApplyFunc(uuid.NewString, func() string {
				return "456"
			})
			defer transactionId.Reset()

			transactionService := transaction.NewTransactionService(mockTransactionRepository, mockUserRepository)
			got, err := transactionService.UpdateTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.UpdateTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.UpdateTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionService_RemoveTransaction(t *testing.T) {
	type args struct {
		ctx context.Context
		req string
	}
	type mockFindTransactionByIDRepository struct {
		res entity.Transaction
		err error
	}
	type mockDeleteTransaction struct {
		err error
	}
	tests := []struct {
		name                              string
		args                              args
		mockFindTransactionByIDRepository *mockFindTransactionByIDRepository
		mockDeleteTransaction             *mockDeleteTransaction
		wantErr                           bool
	}{
		{
			name: "Remove Transaction Success",
			args: args{
				ctx: context.TODO(),
				req: "456",
			},
			mockFindTransactionByIDRepository: &mockFindTransactionByIDRepository{
				res: entity.Transaction{
					TransactionID: "456",
					Name:          "product_test",
					UserID:        "123",
				},
				err: nil,
			},
			mockDeleteTransaction: &mockDeleteTransaction{
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "Error When Find Transaction By ID",
			args: args{
				ctx: context.TODO(),
				req: "4567",
			},
			mockFindTransactionByIDRepository: &mockFindTransactionByIDRepository{
				res: entity.Transaction{},
				err: errors.New("transaction not found"),
			},
			wantErr: true,
		},
		{
			name: "Error When Remove Transaction",
			args: args{
				ctx: context.TODO(),
				req: "456",
			},
			mockFindTransactionByIDRepository: &mockFindTransactionByIDRepository{
				res: entity.Transaction{
					TransactionID: "456",
					Name:          "product_test",
					UserID:        "123",
				},
				err: nil,
			},
			mockDeleteTransaction: &mockDeleteTransaction{
				err: errors.New("error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)

			if tt.mockFindTransactionByIDRepository != nil {
				mockTransactionRepository.On("FindTransactionByID", tt.args.ctx, mock.Anything).Return(tt.mockFindTransactionByIDRepository.res, tt.mockFindTransactionByIDRepository.err)
			}
			if tt.mockDeleteTransaction != nil {
				mockTransactionRepository.On("DeleteTransaction", tt.args.ctx, mock.Anything).Return(tt.mockDeleteTransaction.err)
			}

			transactionService := transaction.NewTransactionService(mockTransactionRepository, mockUserRepository)
			err := transactionService.RemoveTransaction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetTransactionById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
