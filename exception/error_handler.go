package exception

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/vnnyx/golang-dot-api/model/web"
)

func ErrorHandler(err error, ctx echo.Context) {
	if databaseError(err, ctx) {
		return
	}
	if validationError(err, ctx) {
		return
	}
	generalError(err, ctx)
}

func generalError(err error, ctx echo.Context) {
	switch err.Error() {
	case "USER_NOT_FOUND":
		_ = ctx.JSON(http.StatusNotFound, web.WebResponse{
			Code:   http.StatusNotFound,
			Status: web.NOT_FOUND,
			Data:   nil,
			Error: map[string]interface{}{
				"user_id": "NOT_FOUND",
			},
		})
	case "TRANSACTION_NOT_FOUND":
		_ = ctx.JSON(http.StatusNotFound, web.WebResponse{
			Code:   http.StatusNotFound,
			Status: web.NOT_FOUND,
			Data:   nil,
			Error: map[string]interface{}{
				"trasaction_id": "NOT_FOUND",
			},
		})
	case web.UNAUTHORIZATION:
		_ = ctx.JSON(http.StatusUnauthorized, web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: web.UNAUTHORIZATION,
			Data:   nil,
			Error: map[string]interface{}{
				"message": "Unauthorized",
			},
		})
	case "PASSWORD_NOT_MATCH":
		_ = ctx.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: web.BAD_REQUEST,
			Data:   nil,
			Error: map[string]interface{}{
				"password_confirmation": "not match",
			},
		})
	case "code=404, message=Not Found":
		_ = ctx.JSON(http.StatusNotFound, web.WebResponse{
			Code:   http.StatusNotFound,
			Status: web.NOT_FOUND,
			Data:   nil,
			Error: map[string]interface{}{
				"message": "Not Found",
			},
		})
	case "code=405, message=Method Not Allowed":
		_ = ctx.JSON(http.StatusMethodNotAllowed, web.WebResponse{
			Code:   http.StatusMethodNotAllowed,
			Status: web.METHOD_NOT_ALLOWED,
			Data:   nil,
			Error: map[string]interface{}{
				"message": "Method Not Allowed",
			},
		})
	default:
		_ = ctx.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: web.SERVER_ERROR,
			Data:   nil,
			Error: map[string]interface{}{
				"message": "Internal server error",
			},
		})
	}
}

func validationError(err error, ctx echo.Context) bool {
	_, ok := err.(ValidationError)
	if ok {
		var obj interface{}
		_ = json.Unmarshal([]byte(err.Error()), &obj)
		_ = ctx.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: web.BAD_REQUEST,
			Data:   nil,
			Error:  obj,
		})
		return true
	}
	return false
}

func databaseError(err error, ctx echo.Context) bool {
	sqlError, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}
	if sqlError.Number == 1062 && strings.Contains(sqlError.Message, "username") {
		_ = ctx.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: web.BAD_REQUEST,
			Data:   nil,
			Error: map[string]interface{}{
				"username": "must be unique",
			},
		})
	} else if sqlError.Number == 1062 && strings.Contains(sqlError.Message, "email") {
		_ = ctx.JSON(http.StatusBadRequest, web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: web.BAD_REQUEST,
			Data:   nil,
			Error: map[string]interface{}{
				"email": "must be unique",
			},
		})
	} else {
		_ = ctx.JSON(http.StatusInternalServerError, web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: web.SERVER_ERROR,
			Data:   nil,
			Error: map[string]interface{}{
				"message": "Internal server error",
			},
		})
	}
	return true
}
