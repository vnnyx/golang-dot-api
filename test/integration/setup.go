package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vnnyx/golang-dot-api/exception"
	"github.com/vnnyx/golang-dot-api/infrastructure"
	"github.com/vnnyx/golang-dot-api/injector/wire"
	"github.com/vnnyx/golang-dot-api/migration"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/repository/auth"
	"github.com/vnnyx/golang-dot-api/repository/transaction"
	"github.com/vnnyx/golang-dot-api/repository/user"
)

var (
	configuration         = infrastructure.NewConfig(".env.test")
	databases             = infrastructure.NewMySQLDatabase(configuration)
	redis                 = infrastructure.NewRedisClient(".env.test")
	userController        = wire.InitializeUserController(".env.test")
	transactionController = wire.InitializeTransactionController(".env.test")
	authController        = wire.InitializeAuthController(".env.test")
	app                   = testApp()
	userRepository        = user.NewUserRepository(databases)
	transactionRepository = transaction.NewTransactionRepository(databases)
	authRepository        = auth.NewAuthRepository(redis)
	ctx                   = context.TODO()
)

func getAuthorization(payload web.LoginRequest) string {
	requestBody, _ := json.Marshal(payload)

	request := httptest.NewRequest("POST", "/dot-api/login", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	recorder := httptest.NewRecorder()

	app.ServeHTTP(recorder, request)
	response := recorder.Result()

	responseBody, _ := io.ReadAll(response.Body)

	webResponse := web.WebResponse{}
	json.Unmarshal(responseBody, &webResponse)

	jsonData, _ := json.Marshal(webResponse.Data)
	loginResponse := web.LoginResponse{}
	json.Unmarshal(jsonData, &loginResponse)

	return loginResponse.AccessToken
}

func testApp() *echo.Echo {
	migration.Migrate(databases, entity.Transaction{}, entity.User{})
	var app = echo.New()
	app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{DisablePrintStack: true}))
	app.Use(middleware.CORS())
	app.HTTPErrorHandler = exception.ErrorHandler
	userController.Route(app)
	transactionController.Route(app)
	authController.Route(app)
	return app
}
