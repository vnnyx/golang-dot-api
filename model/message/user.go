package message

import (
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
)

var (
	USER_OTP_TOPIC = string("USER.SEND.OTP")
)

type Message struct {
	User entity.User
	OTP  *web.UserEmailVerification
}
