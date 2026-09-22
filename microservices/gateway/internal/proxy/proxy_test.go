package proxy_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aantonioprado/go-architecture/microservices/gateway/internal/proxy"
)

func TestNew_ForwardsRequestAndResponse(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users" {
			t.Errorf("expected path %q, got %q", "/users", r.URL.Path)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer backend.Close()

	p, err := proxy.New(backend.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users", nil)
	rec := httptest.NewRecorder()

	p.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(body) != `{"ok":true}` {
		t.Errorf("expected proxied body, got %q", string(body))
	}
}

func TestNew_InvalidURL(t *testing.T) {
	if _, err := proxy.New("://not-a-url"); err == nil {
		t.Fatal("expected error for an invalid target URL")
	}
}
