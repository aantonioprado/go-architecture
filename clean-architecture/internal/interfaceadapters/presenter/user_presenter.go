package presenter

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/dto"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

// HTTPUserPresenter implements usecases.UserOutputPort. It is the only
// place in this architecture that writes an HTTP response: the interactor
// never touches http.ResponseWriter directly.
type HTTPUserPresenter struct {
	w http.ResponseWriter
}

func NewHTTPUserPresenter(w http.ResponseWriter) *HTTPUserPresenter {
	return &HTTPUserPresenter{w: w}
}

func (p *HTTPUserPresenter) PresentUserCreated(output usecases.UserOutput) {
	p.writeJSON(http.StatusCreated, toUserResponse(output))
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
