package controller

import (
	"github.com/aantonioprado/go-architecture/mvc/internal/dto"
	"github.com/aantonioprado/go-architecture/mvc/internal/response"
	"net/http"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (h *HealthController) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	res := dto.HealthResponse{
		Status: "OK",
	}

	response.JSON(w, http.StatusOK, res)
}
