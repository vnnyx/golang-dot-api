package auth

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model/web"
)

type AuthService interface {
	Login(ctx context.Context, request web.LoginRequest) (response web.LoginResponse, err error)
	Logout(ctx context.Context, accessUUID string) error
}
