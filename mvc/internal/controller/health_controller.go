package controller

import (
	"go-architecture-mvc/internal/dto"
	"go-architecture-mvc/internal/response"
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

	response.Send(w, http.StatusOK, res, "")
}
