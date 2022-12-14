package auth

import (
	"context"
	"errors"

	"github.com/vnnyx/golang-dot-api/infrastructure"
	"github.com/vnnyx/golang-dot-api/model"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/repository/auth"
	"github.com/vnnyx/golang-dot-api/repository/user"
	"github.com/vnnyx/golang-dot-api/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	*infrastructure.Config
	*gorm.DB
	user.UserRepository
	auth.AuthRepository
}

func NewAuthService(config *infrastructure.Config, db *gorm.DB, userRepository user.UserRepository, authRepository auth.AuthRepository) AuthService {
	return &AuthServiceImpl{Config: config, DB: db, UserRepository: userRepository, AuthRepository: authRepository}
}

func (service *AuthServiceImpl) Login(ctx context.Context, request web.LoginRequest) (response web.LoginResponse, err error) {
	user, err := service.UserRepository.FindUserByUsername(ctx, request.Username)
	if err != nil {
		return response, errors.New(web.UNAUTHORIZATION)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return response, errors.New(web.UNAUTHORIZATION)
	}

	td := util.CreateToken(model.JwtPayload{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
	}, service.Config)

	tokenDetails := &model.TokenDetails{
		AccessToken: td.AccessToken,
		AccessUUID:  td.AccessUUID,
		AtExpires:   td.AtExpires,
	}

	err = service.AuthRepository.StoreToken(ctx, *tokenDetails)
	if err != nil {
		return response, err
	}

	response = web.LoginResponse{
		AccessToken: td.AccessToken,
		UserID:      user.UserID,
		Username:    user.Username,
		Email:       user.Email,
	}

	return response, nil

}

func (service *AuthServiceImpl) Logout(ctx context.Context, accessUUID string) error {
	err := service.AuthRepository.DeleteToken(ctx, accessUUID)
	if err != nil {
		return errors.New(web.UNAUTHORIZATION)
	}
	return nil
}
