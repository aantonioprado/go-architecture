package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	user, err := c.service.CreateUser(req.Name, req.Email)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toUserResponse(user))
}

func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.service.ListUsers()
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

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := c.service.GetUser(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, err)
		return
	}

	user, err := c.service.UpdateUser(chi.URLParam(r, "id"), req.Name, req.Email)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toUserResponse(user))
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := c.service.DeleteUser(chi.URLParam(r, "id")); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: err.Error()})
	case errors.Is(err, ErrEmailTaken):
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
}

func toUserResponse(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}
