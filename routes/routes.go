package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/vnnyx/golang-dot-api/handler"
)

func NewRoutes(handler handler.Handler, e *echo.Echo) {
	api := e.Group("/api")
	api.POST("/api/user", handler.UserController.CreateUser)
	api.GET("/api/user/:id", handler.UserController.GetUserById)
	api.GET("/api/user", handler.UserController.GetAllUser)
	api.PUT("/api/user/:id", handler.UserController.UpdateUserProfile)
	api.DELETE("/api/user/:id", handler.UserController.RemoveUser)
}
