package health

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/microservices/gateway/internal/response"
)

// GetHealthCheck reports only the gateway's own liveness, it does not
// aggregate command-service or query-service, to keep this simple.
type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, HealthResponse{Status: "OK"})
}
