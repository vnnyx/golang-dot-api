package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/repository/transaction"
	"github.com/vnnyx/golang-dot-api/repository/user"
	"github.com/vnnyx/golang-dot-api/validation"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	user.UserRepository
	transaction.TransactionRepository
	*gorm.DB
}

func NewUserService(userRepository user.UserRepository, transactionRepository transaction.TransactionRepository, DB *gorm.DB) UserService {
	return &UserServiceImpl{UserRepository: userRepository, TransactionRepository: transactionRepository, DB: DB}
}

func (service *UserServiceImpl) CreateUser(ctx context.Context, request web.UserCreateRequest) (response web.UserResponse, err error) {
	validation.CreateUserValidation(request)

	if request.Password != request.PasswordConfirmation {
		return response, errors.New("PASSWORD_NOT_MATCH")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return response, err
	}

	user := entity.User{
		UserID:    uuid.NewString(),
		Username:  request.Username,
		Email:     request.Email,
		Handphone: request.Handphone,
		Password:  string(password),
	}

	user, err = service.UserRepository.InsertUser(ctx, user)
	if err != nil {
		return response, err
	}

	response = web.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Handphone: user.Handphone,
	}

	return response, nil
}

func (service *UserServiceImpl) GetUserById(ctx context.Context, userId string) (response web.UserResponse, err error) {
	user, err := service.UserRepository.FindUserByID(ctx, userId)
	if err != nil {
		return response, errors.New("USER_NOT_FOUND")
	}

	response = web.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Handphone: user.Handphone,
	}

	return response, nil
}

func (service *UserServiceImpl) GetAllUser(ctx context.Context) (response []web.UserResponse, err error) {
	users, err := service.UserRepository.FindAllUser(ctx)
	if err != nil {
		return response, err
	}

	for _, user := range users {
		response = append(response, web.UserResponse{
			UserID:    user.UserID,
			Username:  user.Username,
			Email:     user.Email,
			Handphone: user.Handphone,
		})
	}

	return response, nil
}

func (service *UserServiceImpl) UpdateUserProfile(ctx context.Context, request web.UserUpdateProfileRequest) (response web.UserResponse, err error) {
	validation.UpdateUserProfileValidation(request)

	user, err := service.UserRepository.FindUserByID(ctx, request.UserID)
	if err != nil {
		return response, errors.New("USER_NOT_FOUND")
	}

	user, err = service.UserRepository.UpdateUser(ctx, entity.User{
		UserID:    user.UserID,
		Username:  request.Username,
		Email:     request.Email,
		Handphone: request.Handphone,
	})

	if err != nil {
		return response, err
	}

	response = web.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Handphone: user.Handphone,
	}

	return response, nil
}

func (service *UserServiceImpl) RemoveUser(ctx context.Context, userId string) error {
	user, err := service.UserRepository.FindUserByID(ctx, userId)
	if err != nil {
		return errors.New("USER_NOT_FOUND")
	}

	tx := service.DB.Begin()
	err = tx.Error
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	err = service.TransactionRepository.DeleteTransactionByUserId(ctx, tx, user.UserID)
	if err != nil {
		return err
	}

	err = service.UserRepository.DeleteUser(ctx, tx, user.UserID)
	if err != nil {
		return err
	}

	return nil
}
