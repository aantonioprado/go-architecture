package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aantonioprado/go-architecture/mvc/internal/dto"
	"github.com/aantonioprado/go-architecture/mvc/internal/model"
	"github.com/aantonioprado/go-architecture/mvc/internal/repository"
	"github.com/aantonioprado/go-architecture/mvc/internal/response"
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
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := model.NewUser(req.Name, req.Email)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Create(user); err != nil {
		if errors.Is(err, repository.ErrEmailTaken) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}

		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, toUserResponse(user))
}

func (h *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
}

func (h *UserController) FindUserByID(w http.ResponseWriter, r *http.Request) {
}

func (h *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
}

func (h *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
}

func toUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}
