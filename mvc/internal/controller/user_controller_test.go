package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/mvc/internal/controller"
	"github.com/aantonioprado/go-architecture/mvc/internal/dto"
	"github.com/aantonioprado/go-architecture/mvc/internal/repository"
)

func newRequestWithID(method, target, id string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestUserController_CreateUser(t *testing.T) {
	ctrl := controller.NewUserController(repository.NewUserRepository())

	body, _ := json.Marshal(dto.CreateUserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ctrl.CreateUser(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var res dto.UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Email != "antonio@antonioeprado.dev" {
		t.Errorf("expected email %q, got %q", "antonio@antonioeprado.dev", res.Email)
	}

	if res.ID == "" {
		t.Error("expected generated ID, got empty string")
	}
}

func TestUserController_CreateUser_MissingName(t *testing.T) {
	ctrl := controller.NewUserController(repository.NewUserRepository())

	body, _ := json.Marshal(dto.CreateUserRequest{
		Email: "antonio@antonioeprado.dev",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ctrl.CreateUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUserController_CreateUser_DuplicateEmail(t *testing.T) {
	repo := repository.NewUserRepository()
	ctrl := controller.NewUserController(repo)

	body, _ := json.Marshal(dto.CreateUserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	ctrl.CreateUser(httptest.NewRecorder(), req)

	req = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	ctrl.CreateUser(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
}
