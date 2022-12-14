package transaction

import "github.com/labstack/echo/v4"

type TransactionController interface {
	Route(e *echo.Echo)
	CreateTransaction(c echo.Context) error
	GetTransactionById(c echo.Context) error
	GetAllTransaction(c echo.Context) error
	GetTransactionByUserId(c echo.Context) error
	UpdateTransaction(c echo.Context) error
	RemoveTransaction(c echo.Context) error
}
