package user_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/user"
)

func newRequestWithID(method, target, id string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestUserController_ListUsers(t *testing.T) {
	repo := user.NewInMemoryUserRepository()
	svc := user.NewReadService(repo)
	ctrl := user.NewUserController(svc)

	if err := svc.Replicate(user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	ctrl.ListUsers(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res []user.UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 user, got %d", len(res))
	}
}

func TestUserController_GetUser(t *testing.T) {
	repo := user.NewInMemoryUserRepository()
	svc := user.NewReadService(repo)
	ctrl := user.NewUserController(svc)

	if err := svc.Replicate(user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := newRequestWithID(http.MethodGet, "/users/1", "1", nil)
	rec := httptest.NewRecorder()

	ctrl.GetUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestUserController_GetUser_NotFound(t *testing.T) {
	repo := user.NewInMemoryUserRepository()
	svc := user.NewReadService(repo)
	ctrl := user.NewUserController(svc)

	req := newRequestWithID(http.MethodGet, "/users/unknown-id", "unknown-id", nil)
	rec := httptest.NewRecorder()

	ctrl.GetUser(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
