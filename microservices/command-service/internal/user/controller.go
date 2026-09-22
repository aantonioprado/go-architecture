package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/response"
)

// UserController is write-only: this service is the command side of the CQRS
// split, reads are served by query-service.
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

	response.JSON(w, http.StatusCreated, toUserResponse(user))
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

	response.JSON(w, http.StatusOK, toUserResponse(user))
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if err := c.service.DeleteUser(chi.URLParam(r, "id")); err != nil {
		writeError(w, err)
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		response.JSON(w, http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
	case errors.Is(err, ErrEmailTaken):
		response.JSON(w, http.StatusConflict, response.ErrorResponse{Error: err.Error()})
	default:
		response.JSON(w, http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
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
