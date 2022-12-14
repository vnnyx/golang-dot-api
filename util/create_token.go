package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/vnnyx/golang-dot-api/exception"
	"github.com/vnnyx/golang-dot-api/infrastructure"
	"github.com/vnnyx/golang-dot-api/model"
)

func CreateToken(request model.JwtPayload, configuration *infrastructure.Config) *model.TokenDetails {
	accessExpired := configuration.JWTMinute

	td := &model.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * time.Duration(accessExpired)).Unix()
	td.AccessUUID = uuid.NewString()

	keyAccess, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(configuration.JWTSecretKey))
	exception.PanicIfNeeded(err)

	now := time.Now().UTC()

	atClaims := jwt.MapClaims{}
	atClaims["id"] = request.UserID
	atClaims["username"] = request.Username
	atClaims["email"] = request.Email
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["exp"] = td.AtExpires
	atClaims["iat"] = now.Unix()
	atClaims["iss"] = "dot-api"
	atClaims["aud"] = "dot-api"

	at := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)
	at.Header["dot-api"] = "jwt"
	td.AccessToken, err = at.SignedString(keyAccess)

	if err != nil {
		exception.PanicIfNeeded(errors.New("AUTHENTICATION_FAILURE"))
	}

	return td
}
