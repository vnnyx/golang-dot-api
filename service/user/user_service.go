package user

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model/web"
)

type UserService interface {
	CreateUser(ctx context.Context, request web.UserCreateRequest) (response web.UserResponse, err error)
	GetUserById(ctx context.Context, userId string) (response web.UserResponse, err error)
	GetAllUser(ctx context.Context) (response []web.UserResponse, err error)
	UpdateUserProfile(ctx context.Context, request web.UserUpdateProfileRequest) (response web.UserResponse, err error)
	RemoveUser(ctx context.Context, userId string) error
}
