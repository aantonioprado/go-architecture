package httpadapter

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/domain"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/ports"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	user, err := h.service.CreateUser(req.Name, req.Email)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		writeError(w, err)
		return
	}

	res := make([]UserResponse, 0, len(users))
	for _, user := range users {
		res = append(res, toUserResponse(user))
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.service.GetUser(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	user, err := h.service.UpdateUser(chi.URLParam(r, "id"), req.Name, req.Email)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteUser(chi.URLParam(r, "id")); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrEmailTaken):
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
}

func toUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}
