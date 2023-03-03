package user

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model/web"
)

type UserService interface {
	CreateUser(ctx context.Context, request web.UserCreateRequest) (response web.UserEmailVerification, err error)
	GetUserById(ctx context.Context, userId string) (response web.UserResponse, err error)
	GetAllUserWithLastTransaction(ctx context.Context, p web.Pagination) (response *web.Pagination, err error)
	GetAllUser(ctx context.Context, p web.Pagination) (response *web.Pagination, err error)
	UpdateUserProfile(ctx context.Context, request web.UserUpdateProfileRequest) (response web.UserResponse, err error)
	RemoveUser(ctx context.Context, userId string) error
	ValidateOTP(ctx context.Context, check web.UserEmailVerification) (response web.UserResponse, err error)
}
