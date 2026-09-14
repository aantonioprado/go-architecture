package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aantonioprado/go-architecture/layered/internal/dto"
	"github.com/aantonioprado/go-architecture/layered/internal/model"
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
	case errors.Is(err, service.ErrNameRequired), errors.Is(err, service.ErrEmailRequired):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrEmailTaken):
		response.Error(w, http.StatusConflict, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, err.Error())
	}
}
