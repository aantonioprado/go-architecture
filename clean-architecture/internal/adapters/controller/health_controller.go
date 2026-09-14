package controller

import (
	"encoding/json"
	"net/http"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/adapters/dto"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(dto.HealthResponse{Status: "OK"})
}
