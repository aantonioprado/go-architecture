package health

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/response"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, HealthResponse{Status: "OK"})
}
