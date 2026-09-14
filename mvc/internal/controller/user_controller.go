package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

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
		writeRepositoryError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, toUserResponse(user))
}

func (h *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.FindAll()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		res = append(res, toUserResponse(user))
	}

	response.JSON(w, http.StatusOK, res)
}

func (h *UserController) FindUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.repo.FindById(id)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toUserResponse(user))
}

func (h *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.repo.FindById(id)
	if err != nil {
		writeRepositoryError(w, err)
		return
	}

	updated, err := user.Update(req.Name, req.Email)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Update(updated); err != nil {
		writeRepositoryError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toUserResponse(updated))
}

func (h *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.repo.Delete(id); err != nil {
		writeRepositoryError(w, err)
		return
	}

	response.NoContent(w)
}

func toUserResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}

func writeRepositoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, repository.ErrEmailTaken):
		response.Error(w, http.StatusConflict, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, err.Error())
	}
}
