package user

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user entity.User) (entity.User, error)
	FindUserByID(ctx context.Context, userId string) (user entity.User, err error)
	FindAllUser(ctx context.Context, p *web.Pagination) (users []*entity.User, err error)
	FindUserByUsername(ctx context.Context, username string) (user entity.User, err error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error
	DeleteAllUser(ctx context.Context) error
	StoreToRedis(ctx context.Context, token web.UserEmailVerification, user entity.User) error
	GetDataToVerify(ctx context.Context, id string) (otp string, user entity.User, err error)
	DeleteCache(ctx context.Context, id string) error
}
