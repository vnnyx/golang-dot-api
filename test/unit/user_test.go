package unit

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	mockTransactionRepository "github.com/vnnyx/golang-dot-api/repository/transaction/mocks"
	mockUserRepository "github.com/vnnyx/golang-dot-api/repository/user/mocks"
	"github.com/vnnyx/golang-dot-api/service/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserService_CreateUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req web.UserCreateRequest
	}
	type mockCreateUserRepository struct {
		res entity.User
		err error
	}
	tests := []struct {
		name                     string
		args                     args
		mockCreateUserRepository *mockCreateUserRepository
		want                     web.UserResponse
		wantErrGeneratePassword  bool
		wantErr                  bool
	}{
		{
			name: "UserService CreateUser Success",
			args: args{
				ctx: context.TODO(),
				req: web.UserCreateRequest{
					Username:             "username_test",
					Email:                "email@test.com",
					Handphone:            "08123456789",
					Password:             "password",
					PasswordConfirmation: "password",
				},
			},
			mockCreateUserRepository: &mockCreateUserRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: nil,
			},
			want: web.UserResponse{
				UserID:    "123",
				Username:  "username_test",
				Email:     "email@test.com",
				Handphone: "08123456789",
			},
			wantErrGeneratePassword: false,
			wantErr:                 false,
		},
		{
			name: "Error When Hashing Password",
			args: args{
				ctx: context.TODO(),
				req: web.UserCreateRequest{
					Username:             "username_test",
					Email:                "email@test.com",
					Handphone:            "08123456789",
					Password:             "password",
					PasswordConfirmation: "password",
				},
			},
			wantErrGeneratePassword: true,
			wantErr:                 true,
		},
		{
			name: "Error When Insert Data To DB",
			args: args{
				ctx: context.TODO(),
				req: web.UserCreateRequest{
					Username:             "username_test",
					Email:                "email@test.com",
					Handphone:            "08123456789",
					Password:             "password",
					PasswordConfirmation: "password",
				},
			},
			mockCreateUserRepository: &mockCreateUserRepository{
				res: entity.User{},
				err: errors.New("error"),
			},
			wantErrGeneratePassword: false,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)
			db, _, err := sqlmock.New()
			require.NoError(t, err)
			DB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			require.NoError(t, err)
			defer db.Close()

			if tt.mockCreateUserRepository != nil {
				mockUserRepository.On("InsertUser", tt.args.ctx, mock.Anything).Return(tt.mockCreateUserRepository.res, tt.mockCreateUserRepository.err)
			}

			if tt.wantErrGeneratePassword {
				password := gomonkey.ApplyFunc(bcrypt.GenerateFromPassword, func(_ []byte, _ int) ([]byte, error) {
					return []byte{}, errors.New("error")
				})
				defer password.Reset()
			}

			userId := gomonkey.ApplyFunc(uuid.NewString, func() string {
				return "123"
			})
			defer userId.Reset()

			userService := user.NewUserService(mockUserRepository, mockTransactionRepository, DB)
			got, err := userService.CreateUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetUserById(t *testing.T) {
	type args struct {
		ctx context.Context
		req string
	}
	type mockFindUserByIdRepository struct {
		res entity.User
		err error
	}
	tests := []struct {
		name                       string
		args                       args
		mockFindUserByIdRepository *mockFindUserByIdRepository
		want                       web.UserResponse
		wantErr                    bool
	}{
		{
			name: "UserService GetUserById Success",
			args: args{
				ctx: context.TODO(),
				req: "123",
			},
			mockFindUserByIdRepository: &mockFindUserByIdRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: nil,
			},
			want: web.UserResponse{
				UserID:    "123",
				Username:  "username_test",
				Email:     "email@test.com",
				Handphone: "08123456789",
			},
			wantErr: false,
		},
		{
			name: "Error Record Not Found",
			args: args{
				ctx: context.TODO(),
				req: "90",
			},
			mockFindUserByIdRepository: &mockFindUserByIdRepository{
				res: entity.User{},
				err: errors.New("record not found"),
			},
			want:    web.UserResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)
			db, _, err := sqlmock.New()
			require.NoError(t, err)
			DB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			require.NoError(t, err)
			defer db.Close()

			if tt.mockFindUserByIdRepository != nil {
				mockUserRepository.On("FindUserByID", tt.args.ctx, mock.Anything).Return(tt.mockFindUserByIdRepository.res, tt.mockFindUserByIdRepository.err)
			}

			userService := user.NewUserService(mockUserRepository, mockTransactionRepository, DB)
			got, err := userService.GetUserById(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetUserById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetUserById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetAllUser(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type mockFindAllUserRepository struct {
		res []entity.User
		err error
	}
	tests := []struct {
		name                      string
		args                      args
		mockFindAllUserRepository *mockFindAllUserRepository
		want                      []web.UserResponse
		wantErr                   bool
	}{
		{
			name: "UserService GetAllUSer Success",
			args: args{
				ctx: context.TODO(),
			},
			mockFindAllUserRepository: &mockFindAllUserRepository{
				res: []entity.User{
					{
						UserID:    "123",
						Username:  "username_test",
						Email:     "email@test.com",
						Handphone: "08123456789",
					},
				},
				err: nil,
			},
			want: []web.UserResponse{
				{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
			},
			wantErr: false,
		},
		{
			name: "Error When Get Data From DB",
			args: args{
				ctx: context.TODO(),
			},
			mockFindAllUserRepository: &mockFindAllUserRepository{
				res: []entity.User{},
				err: errors.New("error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)
			db, _, err := sqlmock.New()
			require.NoError(t, err)
			DB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			require.NoError(t, err)
			defer db.Close()

			if tt.mockFindAllUserRepository != nil {
				mockUserRepository.On("FindAllUser", tt.args.ctx, mock.Anything).Return(tt.mockFindAllUserRepository.res, tt.mockFindAllUserRepository.err)
			}

			userService := user.NewUserService(mockUserRepository, mockTransactionRepository, DB)
			got, err := userService.GetAllUser(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetAllUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetAllUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateUserProfile(t *testing.T) {
	type args struct {
		ctx context.Context
		req web.UserUpdateProfileRequest
	}
	type mockFindUserByIDRepository struct {
		res entity.User
		err error
	}
	type mockUpdateUserRepository struct {
		res entity.User
		err error
	}
	tests := []struct {
		name                       string
		args                       args
		mockFindUserByIDRepository *mockFindUserByIDRepository
		mockUpdateUserRepository   *mockUpdateUserRepository
		want                       web.UserResponse
		wantErr                    bool
	}{
		{
			name: "UserService UpdateUserProfile Success",
			args: args{
				ctx: context.TODO(),
				req: web.UserUpdateProfileRequest{
					UserID:    "123",
					Username:  "username_test_updated",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
			},
			mockFindUserByIDRepository: &mockFindUserByIDRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: nil,
			},
			mockUpdateUserRepository: &mockUpdateUserRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test_updated",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: nil,
			},
			want: web.UserResponse{
				UserID:    "123",
				Username:  "username_test_updated",
				Email:     "email@test.com",
				Handphone: "08123456789",
			},
			wantErr: false,
		},
		{
			name: "Error When Find Record",
			args: args{
				ctx: context.TODO(),
				req: web.UserUpdateProfileRequest{
					UserID:    "123",
					Username:  "username_test_updated",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
			},
			mockFindUserByIDRepository: &mockFindUserByIDRepository{
				res: entity.User{},
				err: errors.New("error"),
			},
			want:    web.UserResponse{},
			wantErr: true,
		},
		{
			name: "Error When Updated Record",
			args: args{
				ctx: context.TODO(),
				req: web.UserUpdateProfileRequest{
					UserID:    "123",
					Username:  "username_test_updated",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
			},
			mockFindUserByIDRepository: &mockFindUserByIDRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: nil,
			},
			mockUpdateUserRepository: &mockUpdateUserRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test_updated",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: errors.New("error"),
			},
			want:    web.UserResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)
			db, _, err := sqlmock.New()
			require.NoError(t, err)
			DB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			require.NoError(t, err)
			defer db.Close()

			if tt.mockFindUserByIDRepository != nil {
				mockUserRepository.On("FindUserByID", tt.args.ctx, mock.Anything).Return(tt.mockFindUserByIDRepository.res, tt.mockFindUserByIDRepository.err)
			}
			if tt.mockUpdateUserRepository != nil {
				mockUserRepository.On("UpdateUser", tt.args.ctx, mock.Anything).Return(tt.mockUpdateUserRepository.res, tt.mockUpdateUserRepository.err)
			}

			userService := user.NewUserService(mockUserRepository, mockTransactionRepository, DB)
			got, err := userService.UpdateUserProfile(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.UpdateUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.UpdateUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_RemoveUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req string
	}
	type mockFindUserByIdRepository struct {
		res entity.User
		err error
	}
	type mockDeleteTransactionByUserIdRespository struct {
		err error
	}
	type mockDeleteUserRepository struct {
		err error
	}
	tests := []struct {
		name                                     string
		args                                     args
		mockFindUserByIdRepository               *mockFindUserByIdRepository
		mockDeleteTransactionByUserIdRespository *mockDeleteTransactionByUserIdRespository
		mockDeleteUserRepository                 *mockDeleteUserRepository
		wantErr                                  bool
	}{
		{
			name: "UserService RemoveUser Success",
			args: args{
				ctx: context.TODO(),
				req: "123",
			},
			mockFindUserByIdRepository: &mockFindUserByIdRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: nil,
			},
			mockDeleteTransactionByUserIdRespository: &mockDeleteTransactionByUserIdRespository{
				err: nil,
			},
			mockDeleteUserRepository: &mockDeleteUserRepository{
				err: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockTransactionRepository := new(mockTransactionRepository.TransactionRepository)
			db, sqlmock, err := sqlmock.New()
			require.NoError(t, err)
			DB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			require.NoError(t, err)
			defer db.Close()

			if tt.mockFindUserByIdRepository != nil {
				mockUserRepository.On("FindUserByID", tt.args.ctx, mock.Anything).Return(tt.mockFindUserByIdRepository.res, tt.mockFindUserByIdRepository.err)
			}

			sqlmock.ExpectBegin()
			if tt.mockDeleteTransactionByUserIdRespository.err != nil || tt.mockDeleteUserRepository.err != nil {
				sqlmock.ExpectRollback()
			} else {
				sqlmock.ExpectCommit()
			}
			if tt.mockDeleteTransactionByUserIdRespository != nil {
				mockTransactionRepository.On("DeleteTransactionByUserId", tt.args.ctx, mock.Anything, tt.args.req).Return(tt.mockDeleteTransactionByUserIdRespository.err)
			}
			if tt.mockDeleteUserRepository != nil {
				mockUserRepository.On("DeleteUser", tt.args.ctx, mock.Anything, tt.args.req).Return(tt.mockDeleteUserRepository.err)
			}

			userService := user.NewUserService(mockUserRepository, mockTransactionRepository, DB)
			err = userService.RemoveUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.RemoveUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
