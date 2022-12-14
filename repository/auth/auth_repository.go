package auth

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model"
)

type AuthRepository interface {
	StoreToken(ctx context.Context, details model.TokenDetails) error
	DeleteToken(ctx context.Context, accessUuid string) error
	GetToken(ctx context.Context, accessUuid string) (access string, err error)
}
