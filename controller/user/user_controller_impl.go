package user

import (
	"net/http"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/labstack/echo/v4"
	"github.com/vnnyx/golang-dot-api/exception"
	authMiddleware "github.com/vnnyx/golang-dot-api/middleware"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/service/user"
)

type UserControllerImpl struct {
	*kafka.Consumer
	user.UserService
	*authMiddleware.AuthMiddleware
}

func NewUserController(userService user.UserService, authMiddleware *authMiddleware.AuthMiddleware, kafka *kafka.Consumer) UserController {
	return &UserControllerImpl{UserService: userService, AuthMiddleware: authMiddleware, Consumer: kafka}
}

func (controller *UserControllerImpl) Route(e *echo.Echo) {
	api := e.Group("/dot-api/user")
	api.POST("", controller.CreateUser)
	api.POST("/verify", controller.VerifyEmail)
	api.GET("/:id", controller.GetUserById)
	api.GET("", controller.GetAllUser)
	api.GET("/transaction", controller.GetAllUserWithLastTransaction)
	api.PUT("/:id", controller.UpdateUserProfile, controller.AuthMiddleware.CheckToken)
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
	var p web.Pagination
	p.Limit, _ = strconv.Atoi(c.QueryParam("limit"))
	p.Page, _ = strconv.Atoi(c.QueryParam("page"))
	p.Sort = c.QueryParam("sort")
	response, err := controller.UserService.GetAllUser(c.Request().Context(), p)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}

func (controller *UserControllerImpl) GetAllUserWithLastTransaction(c echo.Context) error {
	var p web.Pagination
	p.Limit, _ = strconv.Atoi(c.QueryParam("limit"))
	p.Page, _ = strconv.Atoi(c.QueryParam("page"))
	p.Sort = c.QueryParam("sort")
	response, err := controller.UserService.GetAllUserWithLastTransaction(c.Request().Context(), p)
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

func (controller *UserControllerImpl) VerifyEmail(c echo.Context) error {
	var check web.UserEmailVerification
	userId := c.QueryParam("id")
	otp := c.QueryParam("otp")
	check.UserID = userId
	check.OTP, _ = strconv.Atoi(otp)
	response, err := controller.UserService.ValidateOTP(c.Request().Context(), check)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}
