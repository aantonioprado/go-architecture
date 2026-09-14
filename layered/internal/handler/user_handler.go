package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/layered/internal/dto"
	"github.com/aantonioprado/go-architecture/layered/internal/model"
	"github.com/aantonioprado/go-architecture/layered/internal/repository"
	"github.com/aantonioprado/go-architecture/layered/internal/response"
	"github.com/aantonioprado/go-architecture/layered/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Create(req.Name, req.Email)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List()
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

func (h *UserHandler) FindUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.service.GetById(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.Update(id, req.Name, req.Email)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(id); err != nil {
		writeServiceError(w, err)
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

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrNameRequired), errors.Is(err, service.ErrEmailRequired):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrEmailTaken):
		response.Error(w, http.StatusConflict, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, err.Error())
	}
}
