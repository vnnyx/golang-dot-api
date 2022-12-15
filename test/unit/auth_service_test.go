package unit

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vnnyx/golang-dot-api/infrastructure"
	"github.com/vnnyx/golang-dot-api/model"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	mockAuthRepository "github.com/vnnyx/golang-dot-api/repository/auth/mocks"
	mockUserRepository "github.com/vnnyx/golang-dot-api/repository/user/mocks"
	"github.com/vnnyx/golang-dot-api/service/auth"
	"github.com/vnnyx/golang-dot-api/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	config = infrastructure.NewConfig(".env.unit")
)

func TestAuthService_Login(t *testing.T) {
	type args struct {
		ctx context.Context
		req web.LoginRequest
	}
	type mockFindUserByUsernameRepository struct {
		res entity.User
		err error
	}
	type mockStoreTokenRepository struct {
		err error
	}
	tests := []struct {
		name                             string
		args                             args
		token                            *model.TokenDetails
		mockFindUserByUsernameRepository *mockFindUserByUsernameRepository
		mockStoreTokenRepository         *mockStoreTokenRepository
		want                             web.LoginResponse
		wantErrComparePassword           bool
		wantErr                          bool
	}{
		{
			name: "Login Success",
			args: args{
				ctx: context.TODO(),
				req: web.LoginRequest{
					Username: "username_test",
					Password: "password",
				},
			},
			mockFindUserByUsernameRepository: &mockFindUserByUsernameRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: nil,
			},
			mockStoreTokenRepository: &mockStoreTokenRepository{
				err: nil,
			},
			want: web.LoginResponse{
				AccessToken: "access_token",
				UserID:      "123",
				Username:    "username_test",
				Email:       "email@test.com",
			},
			wantErrComparePassword: false,
			wantErr:                false,
		},
		{
			name: "Error When Find User By Username",
			args: args{
				ctx: context.TODO(),
				req: web.LoginRequest{
					Username: "username_test",
					Password: "password",
				},
			},
			mockFindUserByUsernameRepository: &mockFindUserByUsernameRepository{
				res: entity.User{},
				err: errors.New("error"),
			},
			want:    web.LoginResponse{},
			wantErr: true,
		},
		{
			name: "Error When Wrong Input Password",
			args: args{
				ctx: context.TODO(),
				req: web.LoginRequest{
					Username: "username_test",
					Password: "password_wrong",
				},
			},
			mockFindUserByUsernameRepository: &mockFindUserByUsernameRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: nil,
			},
			mockStoreTokenRepository: &mockStoreTokenRepository{
				err: nil,
			},
			want:                   web.LoginResponse{},
			wantErrComparePassword: true,
			wantErr:                true,
		},
		{
			name: "Error When Store Token To Redis",
			args: args{
				ctx: context.TODO(),
				req: web.LoginRequest{
					Username: "username_test",
					Password: "password",
				},
			},
			mockFindUserByUsernameRepository: &mockFindUserByUsernameRepository{
				res: entity.User{
					UserID:    "123",
					Username:  "username_test",
					Email:     "email@test.com",
					Handphone: "08123456789",
				},
				err: nil,
			},
			mockStoreTokenRepository: &mockStoreTokenRepository{
				err: errors.New("error"),
			},
			want:                   web.LoginResponse{},
			wantErrComparePassword: false,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			mockUserRepository := new(mockUserRepository.UserRepository)
			mockAuthRepository := new(mockAuthRepository.AuthRepository)
			db, _, err := sqlmock.New()
			require.NoError(t, err)
			DB, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      db,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{})
			require.NoError(t, err)
			defer db.Close()

			if tt.mockFindUserByUsernameRepository != nil {
				mockUserRepository.On("FindUserByUsername", tt.args.ctx, mock.Anything).Return(tt.mockFindUserByUsernameRepository.res, tt.mockFindUserByUsernameRepository.err)
			}
			if tt.mockStoreTokenRepository != nil {
				mockAuthRepository.On("StoreToken", tt.args.ctx, mock.Anything).Return(tt.mockStoreTokenRepository.err)
			}

			if tt.wantErrComparePassword {
				compare := gomonkey.ApplyFunc(bcrypt.CompareHashAndPassword, func(_ []byte, _ []byte) error {
					return errors.New("error")
				})
				defer compare.Reset()
			} else {
				compare := gomonkey.ApplyFunc(bcrypt.CompareHashAndPassword, func(_ []byte, _ []byte) error {
					return nil
				})
				defer compare.Reset()
			}

			td := gomonkey.ApplyFunc(util.CreateToken, func(_ model.JwtPayload, _ *infrastructure.Config) *model.TokenDetails {
				return &model.TokenDetails{
					AccessToken: "access_token",
					AccessUUID:  "access_uuid",
					AtExpires:   60,
				}
			})
			defer td.Reset()

			authService := auth.NewAuthService(config, DB, mockUserRepository, mockAuthRepository)
			got, err := authService.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
