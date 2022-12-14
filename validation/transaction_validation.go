package validation

import (
	"encoding/json"

	validator "github.com/go-ozzo/ozzo-validation"
	"github.com/vnnyx/golang-dot-api/exception"
	"github.com/vnnyx/golang-dot-api/model/web"
)

func CreateTransactionValidation(request web.TransactionCreateRequest) {
	err := validator.ValidateStruct(&request,
		validator.Field(&request.Name, validator.Required))
	if err != nil {
		b, _ := json.Marshal(err)
		err = exception.ValidationError{
			Message: string(b),
		}
		exception.PanicIfNeeded(err)
	}
}

func UpdateTransactionValidation(request web.TransactionUpdateRequest) {
	err := validator.ValidateStruct(&request,
		validator.Field(&request.Name, validator.Required))
	if err != nil {
		b, _ := json.Marshal(err)
		err = exception.ValidationError{
			Message: string(b),
		}
		exception.PanicIfNeeded(err)
	}
}
