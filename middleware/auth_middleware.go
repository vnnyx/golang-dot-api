package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/vnnyx/golang-dot-api/infrastructure"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/repository/auth"
	"github.com/vnnyx/golang-dot-api/repository/user"
)

type DecodedStructure struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	AccessUUID string `json:"access_uuid"`
}

type AuthMiddleware struct {
	auth.AuthRepository
	user.UserRepository
	ConfigName string
}

func NewAuthMiddleware(authRepository auth.AuthRepository, userRepository user.UserRepository, configName string) *AuthMiddleware {
	return &AuthMiddleware{AuthRepository: authRepository, UserRepository: userRepository, ConfigName: configName}
}

func (middleware *AuthMiddleware) ValidateToken(encodedToken string) (token *jwt.Token, errData error) {
	configuration := infrastructure.NewConfig(middleware.ConfigName)
	jwtPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(configuration.JWTPublicKey))

	if err != nil {
		return token, err
	}

	tokenString := encodedToken
	claims := jwt.MapClaims{}
	token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtPublicKey, nil
	})
	if err != nil {
		return token, err
	}
	if !token.Valid {
		return token, errors.New("invalid token")
	}
	return token, nil
}

func (middleware *AuthMiddleware) DecodeToken(encodedToken string) (decodedResult DecodedStructure, errData error) {
	configuration := infrastructure.NewConfig(middleware.ConfigName)
	jwtPublicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(configuration.JWTPublicKey))
	tokenString := encodedToken
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtPublicKey, nil
	})
	if err != nil {
		return decodedResult, err
	}
	if !token.Valid {
		return decodedResult, errors.New("invalid token")
	}

	jsonbody, err := json.Marshal(claims)
	if err != nil {
		return decodedResult, err
	}

	var obj DecodedStructure
	if err := json.Unmarshal(jsonbody, &obj); err != nil {
		return decodedResult, err
	}

	return obj, nil
}

func (middleware *AuthMiddleware) CheckToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		header := ctx.Request().Header
		tokenSlice := strings.Split(header.Get("Authorization"), "Bearer ")

		var tokenString string
		if len(tokenSlice) == 2 {
			tokenString = tokenSlice[1]
		}

		//validate token
		_, err := middleware.ValidateToken(tokenString)
		if err != nil {
			return errors.New(web.UNAUTHORIZATION)
		}

		//extract data from token
		decodeRes, err := middleware.DecodeToken(tokenString)
		if err != nil {
			return errors.New(web.UNAUTHORIZATION)
		}
		_, err = middleware.UserRepository.FindUserByID(context.Background(), decodeRes.UserID)
		if err != nil {
			return errors.New(web.UNAUTHORIZATION)
		}

		_, err = middleware.AuthRepository.GetToken(context.TODO(), decodeRes.AccessUUID)
		if err != nil {
			return errors.New(web.UNAUTHORIZATION)
		}

		//set global variable
		ctx.Set("currentId", decodeRes.UserID)
		ctx.Set("currentUsername", decodeRes.Username)
		ctx.Set("currentAccessUUID", decodeRes.AccessUUID)

		return next(ctx)
	}
}
