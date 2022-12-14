package user

import (
	"github.com/vnnyx/golang-dot-api/model/entity"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: DB}
}

func (repository *UserRepositoryImpl) InsertUser(user entity.User) (entity.User, error) {
	err := repository.DB.Create(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) FindUserByID(userId string) (user entity.User, err error) {
	err = repository.DB.Where("user_id", userId).First(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) FindUserByUsername(username string) (user entity.User, err error) {
	err = repository.DB.Where("username", username).First(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) FindAllUser() (users []entity.User, err error) {
	err = repository.DB.Find(&users).Error
	return users, err
}

func (repository *UserRepositoryImpl) UpdateUser(user entity.User) (entity.User, error) {
	err := repository.DB.Where("user_id", user.UserID).Updates(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) DeleteUser(tx *gorm.DB, userId string) error {
	return tx.Where("user_id", userId).Delete(&entity.User{}).Error
}
