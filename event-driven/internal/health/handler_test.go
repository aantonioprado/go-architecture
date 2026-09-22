package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aantonioprado/go-architecture/event-driven/internal/health"
)

func TestHandler_GetHealthCheck(t *testing.T) {
	h := health.NewHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.GetHealthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
