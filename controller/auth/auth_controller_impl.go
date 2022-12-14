package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vnnyx/golang-dot-api/exception"
	authMiddleware "github.com/vnnyx/golang-dot-api/middleware"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/service/auth"
)

type AuthControllerImpl struct {
	auth.AuthService
}

func NewAuthController(authService auth.AuthService) AuthController {
	return &AuthControllerImpl{AuthService: authService}
}

func (controller *AuthControllerImpl) Route(e *echo.Echo) {
	api := e.Group("/dot-api")
	api.POST("/login", controller.Login)
	api.POST("/logout", controller.Logout, authMiddleware.CheckToken)
}

func (controller *AuthControllerImpl) Login(c echo.Context) error {
	var request web.LoginRequest
	err := c.Bind(&request)

	response, err := controller.AuthService.Login(c.Request().Context(), request)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}

func (controller *AuthControllerImpl) Logout(c echo.Context) error {
	accessUUID := c.Get("currentAccessUUID")

	err := controller.AuthService.Logout(c.Request().Context(), accessUUID.(string))
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
	})

}
