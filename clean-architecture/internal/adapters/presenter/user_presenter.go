package presenter

import (
	"errors"
	"net/http"
	"time"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/adapters/dto"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/adapters/response"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

type HTTPUserPresenter struct {
	w http.ResponseWriter
}

func NewHTTPUserPresenter(w http.ResponseWriter) *HTTPUserPresenter {
	return &HTTPUserPresenter{w: w}
}

func (p *HTTPUserPresenter) PresentUserCreated(output usecases.UserOutput) {
	response.JSON(p.w, http.StatusCreated, toUserResponse(output))
}

func (p *HTTPUserPresenter) PresentUser(output usecases.UserOutput) {
	response.JSON(p.w, http.StatusOK, toUserResponse(output))
}

func (p *HTTPUserPresenter) PresentUserList(output usecases.ListUsersOutput) {
	res := make([]dto.UserResponse, 0, len(output.Users))
	for _, user := range output.Users {
		res = append(res, toUserResponse(user))
	}

	response.JSON(p.w, http.StatusOK, res)
}

func (p *HTTPUserPresenter) PresentUserDeleted() {
	response.JSON(p.w, http.StatusNoContent, nil)
}

func (p *HTTPUserPresenter) PresentError(err error) {
	switch {
	case errors.Is(err, usecases.ErrUserNotFound):
		response.JSON(p.w, http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, usecases.ErrEmailTaken):
		response.JSON(p.w, http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
	default:
		response.JSON(p.w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	}
}

func toUserResponse(output usecases.UserOutput) dto.UserResponse {
	return dto.UserResponse{
		ID:        output.ID,
		Name:      output.Name,
		Email:     output.Email,
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
	}
}
