package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vnnyx/golang-dot-api/exception"
	"github.com/vnnyx/golang-dot-api/infrastructure"
	"github.com/vnnyx/golang-dot-api/injector/wire"
	"github.com/vnnyx/golang-dot-api/migration"
	"github.com/vnnyx/golang-dot-api/model/entity"
)

func main() {
	configuration := infrastructure.NewConfig(".env")
	databases := infrastructure.NewMySQLDatabase(configuration)
	migration.Migrate(databases, entity.Transaction{}, entity.User{})

	userController := wire.InitializeUserController(".env")
	transactionController := wire.InitializeTransactionController(".env")
	authController := wire.InitializeAuthController(".env")

	app := echo.New()
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisablePrintStack: true}))
	app.Use(middleware.CORS())
	app.HTTPErrorHandler = exception.ErrorHandler
	userController.Route(app)
	transactionController.Route(app)
	authController.Route(app)
	err := app.Start(fmt.Sprintf(":%v", configuration.AppPort))
	exception.PanicIfNeeded(err)
}
