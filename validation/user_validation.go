package validation

import (
	"encoding/json"

	validator "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/vnnyx/golang-dot-api/exception"
	"github.com/vnnyx/golang-dot-api/model/web"
)

func CreateUserValidation(request web.UserCreateRequest) {
	err := validator.ValidateStruct(&request,
		validator.Field(&request.Username, validator.Required),
		validator.Field(&request.Email, validator.Required, is.Email),
		validator.Field(&request.Handphone, validator.Required, is.Digit),
		validator.Field(&request.Password, validator.Required),
		validator.Field(&request.PasswordConfirmation, validator.Required))
	if err != nil {
		b, _ := json.Marshal(err)
		err = exception.ValidationError{
			Message: string(b),
		}
		exception.PanicIfNeeded(err)
	}
}

func UpdateUserProfileValidation(request web.UserUpdateProfileRequest) {
	err := validator.ValidateStruct(&request,
		validator.Field(&request.Username, validator.Required),
		validator.Field(&request.Email, validator.Required, is.Email),
		validator.Field(&request.Handphone, validator.Required, is.Digit))
	if err != nil {
		b, _ := json.Marshal(err)
		err = exception.ValidationError{
			Message: string(b),
		}
		exception.PanicIfNeeded(err)
	}
}
