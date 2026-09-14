package handler

import (
	"net/http"

	"github.com/aantonioprado/go-architecture/layered/internal/dto"
	"github.com/aantonioprado/go-architecture/layered/internal/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	res := dto.HealthResponse{
		Status: "OK",
	}

	response.JSON(w, http.StatusOK, res)
}
