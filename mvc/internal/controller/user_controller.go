package controller

import (
	"go-architecture-mvc/internal/repository"
	"net/http"
)

type UserController struct {
	repo *repository.UserRepository
}

func NewUserController(repo *repository.UserRepository) *UserController {
	return &UserController{
		repo: repo,
	}
}

func (h *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
}

func (h *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
}

func (h *UserController) FindUserByID(w http.ResponseWriter, r *http.Request) {
}

func (h *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
}

func (h *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
}
