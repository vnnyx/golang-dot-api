package user

import (
	"context"

	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/util"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: DB}
}

func (repository *UserRepositoryImpl) InsertUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := repository.DB.WithContext(ctx).Create(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) FindUserByID(ctx context.Context, userId string) (user entity.User, err error) {
	err = repository.DB.WithContext(ctx).Where("id", userId).First(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) FindUserByUsername(ctx context.Context, username string) (user entity.User, err error) {
	err = repository.DB.WithContext(ctx).Where("username", username).First(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) FindAllUser(ctx context.Context, p *web.Pagination) (users []entity.User, err error) {
	if p == nil {
		err = repository.DB.WithContext(ctx).Find(&users).Error
	} else {
		err = repository.DB.WithContext(ctx).Scopes(util.Paginate(&users, p, repository.DB)).Find(&users).Error
	}
	return users, err
}

func (repository *UserRepositoryImpl) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := repository.DB.WithContext(ctx).Where("id", user.UserID).Updates(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error {
	return tx.WithContext(ctx).Where("id", userId).Delete(&entity.User{}).Error
}

func (repository *UserRepositoryImpl) DeleteAllUser(ctx context.Context) error {
	return repository.DB.WithContext(ctx).Exec("DELETE FROM users").Error
}
