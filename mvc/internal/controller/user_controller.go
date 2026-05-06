package controller

import (
	"net/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (h *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
}
