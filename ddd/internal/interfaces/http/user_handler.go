package httpinterface

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	appuser "github.com/aantonioprado/go-architecture/ddd/internal/application/user"
	"github.com/aantonioprado/go-architecture/ddd/internal/domain/user"
)

type UserHandler struct {
	service appuser.Service
}

func NewUserHandler(service appuser.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	u, err := h.service.CreateUser(req.Name, req.Email)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponse(u))
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		writeError(w, err)
		return
	}

	res := make([]UserResponse, 0, len(users))
	for _, u := range users {
		res = append(res, toUserResponse(u))
	}

	writeJSON(w, http.StatusOK, res)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	u, err := h.service.GetUser(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(u))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	u, err := h.service.UpdateUser(chi.URLParam(r, "id"), req.Name, req.Email)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(u))
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
	case errors.Is(err, user.ErrUserNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, user.ErrEmailTaken):
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
}

func toUserResponse(u *user.User) UserResponse {
	return UserResponse{
		ID:        u.ID(),
		Name:      u.Name(),
		Email:     u.Email().String(),
		CreatedAt: u.CreatedAt().Format(time.RFC3339),
	}
}
