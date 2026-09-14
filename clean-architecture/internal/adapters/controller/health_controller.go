package controller

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/adapters/dto"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/adapters/response"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, dto.HealthResponse{Status: "OK"})
}
