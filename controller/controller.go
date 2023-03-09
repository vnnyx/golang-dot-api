package controller

import (
	"github.com/vnnyx/golang-dot-api/controller/auth"
	"github.com/vnnyx/golang-dot-api/controller/transaction"
	"github.com/vnnyx/golang-dot-api/controller/user"
)

type Controller struct {
	transaction.TransactionController
	user.UserController
	auth.AuthController
}

func NewController(transactionController transaction.TransactionController, userController user.UserController, authController auth.AuthController) *Controller {
	return &Controller{TransactionController: transactionController, UserController: userController, AuthController: authController}
}
