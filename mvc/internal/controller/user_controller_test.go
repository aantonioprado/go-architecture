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

func TestUserController_ListUsers(t *testing.T) {
	repo := repository.NewUserRepository()
	ctrl := controller.NewUserController(repo)

	body, _ := json.Marshal(dto.CreateUserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})
	ctrl.CreateUser(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	ctrl.ListUsers(rec, req)

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

func TestUserController_FindUserByID(t *testing.T) {
	repo := repository.NewUserRepository()
	ctrl := controller.NewUserController(repo)

	body, _ := json.Marshal(dto.CreateUserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})
	createRec := httptest.NewRecorder()
	ctrl.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	var created dto.UserResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	req := newRequestWithID(http.MethodGet, "/users/"+created.ID, created.ID, nil)
	rec := httptest.NewRecorder()

	ctrl.FindUserByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestUserController_FindUserByID_NotFound(t *testing.T) {
	ctrl := controller.NewUserController(repository.NewUserRepository())

	req := newRequestWithID(http.MethodGet, "/users/unknown-id", "unknown-id", nil)
	rec := httptest.NewRecorder()

	ctrl.FindUserByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUserController_DeleteUser(t *testing.T) {
	repo := repository.NewUserRepository()
	ctrl := controller.NewUserController(repo)

	createBody, _ := json.Marshal(dto.CreateUserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})
	createRec := httptest.NewRecorder()
	ctrl.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(createBody)))

	var created dto.UserResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	req := newRequestWithID(http.MethodDelete, "/users/"+created.ID, created.ID, nil)
	rec := httptest.NewRecorder()

	ctrl.DeleteUser(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestUserController_DeleteUser_NotFound(t *testing.T) {
	ctrl := controller.NewUserController(repository.NewUserRepository())

	req := newRequestWithID(http.MethodDelete, "/users/unknown-id", "unknown-id", nil)
	rec := httptest.NewRecorder()

	ctrl.DeleteUser(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUserController_UpdateUser(t *testing.T) {
	repo := repository.NewUserRepository()
	ctrl := controller.NewUserController(repo)

	createBody, _ := json.Marshal(dto.CreateUserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})
	createRec := httptest.NewRecorder()
	ctrl.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(createBody)))

	var created dto.UserResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	updateBody, _ := json.Marshal(dto.UpdateUserRequest{
		Name:  "Antônio Elias Prado",
		Email: "antonio@antonioeprado.dev",
	})
	req := newRequestWithID(http.MethodPut, "/users/"+created.ID, created.ID, bytes.NewReader(updateBody))
	rec := httptest.NewRecorder()

	ctrl.UpdateUser(rec, req)

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

func TestUserController_UpdateUser_NotFound(t *testing.T) {
	ctrl := controller.NewUserController(repository.NewUserRepository())

	body, _ := json.Marshal(dto.UpdateUserRequest{
		Name:  "Antônio Elias Prado",
		Email: "antonio@antonioeprado.dev",
	})
	req := newRequestWithID(http.MethodPut, "/users/unknown-id", "unknown-id", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ctrl.UpdateUser(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
