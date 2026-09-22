package httpinterface_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	appuser "github.com/aantonioprado/go-architecture/ddd/internal/application/user"
	"github.com/aantonioprado/go-architecture/ddd/internal/infrastructure/persistence/memory"
	httpinterface "github.com/aantonioprado/go-architecture/ddd/internal/interfaces/http"
)

func newUserHandler() *httpinterface.UserHandler {
	repo := memory.NewUserRepository()
	svc := appuser.NewService(repo)

	return httpinterface.NewUserHandler(svc)
}

func newRequestWithID(method, target, id string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestUserHandler_CreateUser(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(httpinterface.UserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateUser(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var res httpinterface.UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Email != "antonio@antonioeprado.dev" {
		t.Errorf("expected email %q, got %q", "antonio@antonioeprado.dev", res.Email)
	}
}

func TestUserHandler_CreateUser_MissingName(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(httpinterface.UserRequest{
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

	body, _ := json.Marshal(httpinterface.UserRequest{
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

	body, _ := json.Marshal(httpinterface.UserRequest{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"})
	h.CreateUser(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	h.ListUsers(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res []httpinterface.UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 user, got %d", len(res))
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(httpinterface.UserRequest{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"})
	createRec := httptest.NewRecorder()
	h.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	var created httpinterface.UserResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	req := newRequestWithID(http.MethodGet, "/users/"+created.ID, created.ID, nil)
	rec := httptest.NewRecorder()

	h.GetUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	h := newUserHandler()

	req := newRequestWithID(http.MethodGet, "/users/unknown-id", "unknown-id", nil)
	rec := httptest.NewRecorder()

	h.GetUser(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	h := newUserHandler()

	createBody, _ := json.Marshal(httpinterface.UserRequest{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"})
	createRec := httptest.NewRecorder()
	h.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(createBody)))

	var created httpinterface.UserResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	req := newRequestWithID(http.MethodDelete, "/users/"+created.ID, created.ID, nil)
	rec := httptest.NewRecorder()

	h.DeleteUser(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestUserHandler_DeleteUser_NotFound(t *testing.T) {
	h := newUserHandler()

	req := newRequestWithID(http.MethodDelete, "/users/unknown-id", "unknown-id", nil)
	rec := httptest.NewRecorder()

	h.DeleteUser(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	h := newUserHandler()

	createBody, _ := json.Marshal(httpinterface.UserRequest{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"})
	createRec := httptest.NewRecorder()
	h.CreateUser(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(createBody)))

	var created httpinterface.UserResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	updateBody, _ := json.Marshal(httpinterface.UserRequest{Name: "Antônio Elias Prado", Email: "antonio@antonioeprado.dev"})
	req := newRequestWithID(http.MethodPut, "/users/"+created.ID, created.ID, bytes.NewReader(updateBody))
	rec := httptest.NewRecorder()

	h.UpdateUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var updated httpinterface.UserResponse
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if updated.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}
}

func TestUserHandler_UpdateUser_NotFound(t *testing.T) {
	h := newUserHandler()

	body, _ := json.Marshal(httpinterface.UserRequest{Name: "Antônio Elias Prado", Email: "antonio@antonioeprado.dev"})
	req := newRequestWithID(http.MethodPut, "/users/unknown-id", "unknown-id", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.UpdateUser(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
