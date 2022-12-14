package auth

import "github.com/labstack/echo/v4"

type AuthController interface {
	Route(e *echo.Echo)
	Login(c echo.Context) error
	Logout(c echo.Context) error
}
