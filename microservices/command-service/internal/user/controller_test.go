package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/aantonioprado/go-architecture/microservices/command-service/internal/user"
)

func newUserController() *user.UserController {
	repo := user.NewInMemoryUserRepository()
	svc := user.NewUserService(repo, &fakeReplicator{})

	return user.NewUserController(svc)
}

func newRequestWithID(method, target, id string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestUserController_CreateUser(t *testing.T) {
	ctrl := newUserController()

	body, _ := json.Marshal(user.UserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ctrl.CreateUser(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var res user.UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Email != "antonio@antonioeprado.dev" {
		t.Errorf("expected email %q, got %q", "antonio@antonioeprado.dev", res.Email)
	}
}

func TestUserController_CreateUser_MissingName(t *testing.T) {
	ctrl := newUserController()

	body, _ := json.Marshal(user.UserRequest{
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
	ctrl := newUserController()

	body, _ := json.Marshal(user.UserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})

	ctrl.CreateUser(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	rec := httptest.NewRecorder()
	ctrl.CreateUser(rec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
}

func TestUserController_UpdateUser(t *testing.T) {
	ctrl := newUserController()

	createBody, _ := json.Marshal(user.UserRequest{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"})
	createRec := httptest.NewRecorder()
	ctrl.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(createBody)))

	var created user.UserResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	updateBody, _ := json.Marshal(user.UserRequest{Name: "Antônio Elias Prado", Email: "antonio@antonioeprado.dev"})
	req := newRequestWithID(http.MethodPut, "/users/"+created.ID, created.ID, bytes.NewReader(updateBody))
	rec := httptest.NewRecorder()

	ctrl.UpdateUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var updated user.UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if updated.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}
}

func TestUserController_UpdateUser_NotFound(t *testing.T) {
	ctrl := newUserController()

	body, _ := json.Marshal(user.UserRequest{Name: "Antônio Elias Prado", Email: "antonio@antonioeprado.dev"})
	req := newRequestWithID(http.MethodPut, "/users/unknown-id", "unknown-id", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ctrl.UpdateUser(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUserController_DeleteUser(t *testing.T) {
	ctrl := newUserController()

	createBody, _ := json.Marshal(user.UserRequest{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"})
	createRec := httptest.NewRecorder()
	ctrl.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(createBody)))

	var created user.UserResponse
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
	ctrl := newUserController()

	req := newRequestWithID(http.MethodDelete, "/users/unknown-id", "unknown-id", nil)
	rec := httptest.NewRecorder()

	ctrl.DeleteUser(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
