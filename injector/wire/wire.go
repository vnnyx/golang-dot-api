//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	authController "github.com/vnnyx/golang-dot-api/controller/auth"
	transactionController "github.com/vnnyx/golang-dot-api/controller/transaction"
	userController "github.com/vnnyx/golang-dot-api/controller/user"
	"github.com/vnnyx/golang-dot-api/infrastructure"
	authRepository "github.com/vnnyx/golang-dot-api/repository/auth"
	transactionRepository "github.com/vnnyx/golang-dot-api/repository/transaction"
	userRepository "github.com/vnnyx/golang-dot-api/repository/user"
	authService "github.com/vnnyx/golang-dot-api/service/auth"
	transactionService "github.com/vnnyx/golang-dot-api/service/transaction"
	userService "github.com/vnnyx/golang-dot-api/service/user"
)

func InitializeUserController() userController.UserController {
	wire.Build(
		infrastructure.NewConfig,
		infrastructure.NewMySQLDatabase,
		transactionRepository.NewTransactionRepository,
		userRepository.NewUserRepository,
		userService.NewUserService,
		userController.NewUserController,
	)
	return nil
}

func InitializeTransactionController() transactionController.TransactionController {
	wire.Build(
		infrastructure.NewConfig,
		infrastructure.NewMySQLDatabase,
		transactionRepository.NewTransactionRepository,
		userRepository.NewUserRepository,
		transactionService.NewTransactionService,
		transactionController.NewTransactionController,
	)
	return nil
}

func InitializeAuthController() authController.AuthController {
	wire.Build(
		infrastructure.NewConfig,
		infrastructure.NewMySQLDatabase,
		infrastructure.NewRedisClient,
		userRepository.NewUserRepository,
		authRepository.NewAuthRepository,
		authService.NewAuthService,
		authController.NewAuthController,
	)
	return nil
}
