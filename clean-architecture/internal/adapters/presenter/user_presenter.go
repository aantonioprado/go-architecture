package presenter

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/adapters/dto"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

type HTTPUserPresenter struct {
	w http.ResponseWriter
}

func NewHTTPUserPresenter(w http.ResponseWriter) *HTTPUserPresenter {
	return &HTTPUserPresenter{w: w}
}

func (p *HTTPUserPresenter) PresentUserCreated(output usecases.UserOutput) {
	p.writeJSON(http.StatusCreated, toUserResponse(output))
}

func (p *HTTPUserPresenter) PresentUser(output usecases.UserOutput) {
	p.writeJSON(http.StatusOK, toUserResponse(output))
}

func (p *HTTPUserPresenter) PresentUserList(output usecases.ListUsersOutput) {
	res := make([]dto.UserResponse, 0, len(output.Users))
	for _, user := range output.Users {
		res = append(res, toUserResponse(user))
	}

	p.writeJSON(http.StatusOK, res)
}

func (p *HTTPUserPresenter) PresentUserDeleted() {
	p.w.WriteHeader(http.StatusNoContent)
}

func (p *HTTPUserPresenter) PresentError(err error) {
	switch {
	case errors.Is(err, usecases.ErrUserNotFound):
		p.writeJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
	case errors.Is(err, usecases.ErrEmailTaken):
		p.writeJSON(http.StatusConflict, dto.ErrorResponse{Error: err.Error()})
	default:
		p.writeJSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
	}
}

func (p *HTTPUserPresenter) writeJSON(status int, data any) {
	p.w.Header().Set("Content-Type", "application/json")
	p.w.WriteHeader(status)

	if data == nil {
		return
	}

	_ = json.NewEncoder(p.w).Encode(data)
}

func toUserResponse(output usecases.UserOutput) dto.UserResponse {
	return dto.UserResponse{
		ID:        output.ID,
		Name:      output.Name,
		Email:     output.Email,
		CreatedAt: output.CreatedAt.Format(time.RFC3339),
	}
}
