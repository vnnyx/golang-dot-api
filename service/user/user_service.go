package user

import (
	"github.com/vnnyx/golang-dot-api/model/web"
)

type UserService interface {
	CreateUser(request web.UserCreateRequest) (response web.UserResponse, err error)
	GetUserById(userId string) (response web.UserResponse, err error)
	GetAllUser() (response []web.UserResponse, err error)
	UpdateUserProfile(request web.UserUpdateProfileRequest) (response web.UserResponse, err error)
	RemoveUser(userId string) error
}
