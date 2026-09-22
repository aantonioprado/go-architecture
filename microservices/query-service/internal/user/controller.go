package user

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/response"
)

// UserController is read-only: this service is the query side of the CQRS
// split, writes only ever arrive here as replication from command-service
// (see InternalController).
type UserController struct {
	service *ReadService
}

func NewUserController(service *ReadService) *UserController {
	return &UserController{service: service}
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

	response.JSON(w, http.StatusOK, res)
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := c.service.GetUser(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toUserResponse(user))
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		response.JSON(w, http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
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
