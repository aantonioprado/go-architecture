package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/layered/internal/dto"
	"github.com/aantonioprado/go-architecture/layered/internal/handler"
	"github.com/aantonioprado/go-architecture/layered/internal/repository"
	"github.com/aantonioprado/go-architecture/layered/internal/service"
)

func newRequestWithID(method, target, id string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func newUserHandler() *handler.UserHandler {
	return handler.NewUserHandler(service.NewUserService(repository.NewUserRepository()))
}

func TestUserHandler_CreateUser(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(dto.CreateUserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateUser(rec, req)

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
}

func TestUserHandler_CreateUser_MissingName(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(dto.CreateUserRequest{
		Email: "antonio@antonioeprado.dev",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUserHandler_CreateUser_DuplicateEmail(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(dto.CreateUserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})

	h.CreateUser(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	rec := httptest.NewRecorder()
	h.CreateUser(rec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
}

func TestUserHandler_ListUsers(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(dto.CreateUserRequest{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"})
	h.CreateUser(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	h.ListUsers(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res []dto.UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 user, got %d", len(res))
	}
}

func TestUserHandler_FindUserById(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(dto.CreateUserRequest{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"})
	createRec := httptest.NewRecorder()
	h.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	var created dto.UserResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	req := newRequestWithID(http.MethodGet, "/users/"+created.ID, created.ID, nil)
	rec := httptest.NewRecorder()

	h.FindUserById(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestUserHandler_FindUserById_NotFound(t *testing.T) {
	h := newUserHandler()

	req := newRequestWithID(http.MethodGet, "/users/unknown-id", "unknown-id", nil)
	rec := httptest.NewRecorder()

	h.FindUserById(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	h := newUserHandler()

	createBody, _ := json.Marshal(dto.CreateUserRequest{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"})
	createRec := httptest.NewRecorder()
	h.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(createBody)))

	var created dto.UserResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	updateBody, _ := json.Marshal(dto.UpdateUserRequest{Name: "Antônio Elias Prado", Email: "antonio@antonioeprado.dev"})
	req := newRequestWithID(http.MethodPut, "/users/"+created.ID, created.ID, bytes.NewReader(updateBody))
	rec := httptest.NewRecorder()

	h.UpdateUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var updated dto.UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if updated.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}
}

func TestUserHandler_UpdateUser_NotFound(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(dto.UpdateUserRequest{Name: "Antônio Elias Prado", Email: "antonio@antonioeprado.dev"})
	req := newRequestWithID(http.MethodPut, "/users/unknown-id", "unknown-id", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateUser(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
