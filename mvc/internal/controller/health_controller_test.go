package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-architecture-mvc/internal/controller"
)

func TestHealthController_GetHealthCheck(t *testing.T) {
	ctrl := controller.NewHealthController()

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	ctrl.GetHealthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}
