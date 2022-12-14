package user

import (
	"github.com/vnnyx/golang-dot-api/model/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(user entity.User) (entity.User, error)
	FindUserByID(userId string) (user entity.User, err error)
	FindAllUser() (users []entity.User, err error)
	FindUserByUsername(username string) (user entity.User, err error)
	UpdateUser(user entity.User) (entity.User, error)
	DeleteUser(tx *gorm.DB, userId string) error
}
