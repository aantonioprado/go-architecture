package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aantonioprado/go-architecture/layered/internal/handler"
)

func TestHealthHandler_GetHealthCheck(t *testing.T) {
	h := handler.NewHealthHandler()

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	h.GetHealthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}
