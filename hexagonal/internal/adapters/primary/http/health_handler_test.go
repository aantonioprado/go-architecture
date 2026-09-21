package httpadapter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpadapter "github.com/aantonioprado/go-architecture/hexagonal/internal/adapters/primary/http"
)

func TestHealthHandler_GetHealthCheck(t *testing.T) {
	h := httpadapter.NewHealthHandler()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	h.GetHealthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
