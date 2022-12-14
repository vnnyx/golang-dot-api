package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vnnyx/golang-dot-api/exception"
	authMiddleware "github.com/vnnyx/golang-dot-api/middleware"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/service/user"
)

type UserControllerImpl struct {
	user.UserService
}

func NewUserController(userService user.UserService) UserController {
	return &UserControllerImpl{UserService: userService}
}

func (controller *UserControllerImpl) Route(e *echo.Echo) {
	api := e.Group("/dot-api/user")
	api.POST("", controller.CreateUser)
	api.GET("/:id", controller.GetUserById)
	api.GET("", controller.GetAllUser)
	api.PUT("/:id", controller.UpdateUserProfile, authMiddleware.CheckToken)
	api.DELETE("/:id", controller.RemoveUser)
}

func (controller *UserControllerImpl) CreateUser(c echo.Context) error {
	var request web.UserCreateRequest
	err := c.Bind(&request)
	exception.PanicIfNeeded(err)

	response, err := controller.UserService.CreateUser(c.Request().Context(), request)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusCreated, web.WebResponse{
		Code:   http.StatusCreated,
		Status: web.CREATED,
		Data:   response,
	})
}

func (controller *UserControllerImpl) GetUserById(c echo.Context) error {
	userId := c.Param("id")

	response, err := controller.UserService.GetUserById(c.Request().Context(), userId)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}

func (controller *UserControllerImpl) GetAllUser(c echo.Context) error {
	response, err := controller.UserService.GetAllUser(c.Request().Context())
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}

func (controller *UserControllerImpl) UpdateUserProfile(c echo.Context) error {
	var request web.UserUpdateProfileRequest
	err := c.Bind(&request)
	exception.PanicIfNeeded(err)

	request.UserID = c.Param("id")
	response, err := controller.UserService.UpdateUserProfile(c.Request().Context(), request)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}

func (controller *UserControllerImpl) RemoveUser(c echo.Context) error {
	userId := c.Param("id")

	err := controller.UserService.RemoveUser(c.Request().Context(), userId)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
	})
}
