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

	controller, err := wire.InitializeController(".env")
	if err != nil {
		panic(err)
	}
	app := echo.New()
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisablePrintStack: true}))
	app.Use(middleware.CORS())

	app.HTTPErrorHandler = exception.ErrorHandler
	controller.AuthController.Route(app)
	controller.TransactionController.Route(app)
	controller.UserController.Route(app)
	err = app.Start(fmt.Sprintf(":%v", configuration.AppPort))
	exception.PanicIfNeeded(err)
}
