package handler

import "github.com/vnnyx/golang-dot-api/controller/user"

type Handler struct {
	user.UserController
}

func NewHandler(userController user.UserController) *Handler {
	return &Handler{UserController: userController}
}
