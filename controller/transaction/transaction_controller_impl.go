package transaction

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vnnyx/golang-dot-api/exception"
	authMiddleware "github.com/vnnyx/golang-dot-api/middleware"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/service/transaction"
)

type TransactionControllerImpl struct {
	transaction.TransactionService
	*authMiddleware.AuthMiddleware
}

func NewTransactionController(transactionService transaction.TransactionService, authMiddleware *authMiddleware.AuthMiddleware) TransactionController {
	return &TransactionControllerImpl{TransactionService: transactionService, AuthMiddleware: authMiddleware}
}

func (controller *TransactionControllerImpl) Route(e *echo.Echo) {
	api := e.Group("/dot-api/transaction", controller.AuthMiddleware.CheckToken)
	api.POST("", controller.CreateTransaction)
	api.GET("/:id", controller.GetTransactionById)
	api.GET("", controller.GetAllTransaction)
	api.GET("/user", controller.GetTransactionByUserId)
	api.PATCH("/:id", controller.UpdateTransaction)
	api.DELETE("/:id", controller.RemoveTransaction)
}

func (controller *TransactionControllerImpl) CreateTransaction(c echo.Context) error {
	var request web.TransactionCreateRequest
	err := c.Bind(&request)
	exception.PanicIfNeeded(err)

	request.UserID = c.QueryParam("user_id")
	response, err := controller.TransactionService.CreateTransaction(c.Request().Context(), request)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusCreated, web.WebResponse{
		Code:   http.StatusCreated,
		Status: web.CREATED,
		Data:   response,
	})
}

func (controller *TransactionControllerImpl) GetTransactionById(c echo.Context) error {
	transactionId := c.Param("id")

	response, err := controller.TransactionService.GetTransactionById(c.Request().Context(), transactionId)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}

func (controller *TransactionControllerImpl) GetAllTransaction(c echo.Context) error {
	response, err := controller.TransactionService.GetAllTransaction(c.Request().Context())
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}

func (controller *TransactionControllerImpl) GetTransactionByUserId(c echo.Context) error {
	userId := c.QueryParam("user_id")

	response, err := controller.TransactionService.GetTransactionByUserId(c.Request().Context(), userId)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}

func (controller *TransactionControllerImpl) UpdateTransaction(c echo.Context) error {
	var request web.TransactionUpdateRequest
	err := c.Bind(&request)
	exception.PanicIfNeeded(err)

	request.TransactionID = c.Param("id")
	response, err := controller.TransactionService.UpdateTransaction(c.Request().Context(), request)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
		Data:   response,
	})
}

func (controller *TransactionControllerImpl) RemoveTransaction(c echo.Context) error {
	transactionId := c.Param("id")

	err := controller.TransactionService.RemoveTransaction(c.Request().Context(), transactionId)
	exception.PanicIfNeeded(err)

	return c.JSON(http.StatusOK, web.WebResponse{
		Code:   http.StatusOK,
		Status: web.OK,
	})
}
