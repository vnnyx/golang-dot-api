package user

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user entity.User) (entity.User, error)
	FindUserByID(ctx context.Context, userId string) (user entity.User, err error)
	FindAllUser(ctx context.Context) (users []entity.User, err error)
	FindUserByUsername(ctx context.Context, username string) (user entity.User, err error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error
	DeleteAllUser(ctx context.Context) error
}
