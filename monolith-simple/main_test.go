package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	router := newRouter(newUserStore())

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestCreateUser(t *testing.T) {
	router := newRouter(newUserStore())

	body, _ := json.Marshal(createUserRequest{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var res userResponse
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

func TestCreateUser_MissingName(t *testing.T) {
	router := newRouter(newUserStore())

	body, _ := json.Marshal(createUserRequest{
		Email: "semnome@antonioeprado.dev",
	})

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	router := newRouter(newUserStore())

	body, _ := json.Marshal(createUserRequest{
		Name:  "Antônio Prado",
		Email: "duplicado@antonioeprado.dev",
	})

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}
}

func TestListUsers(t *testing.T) {
	router := newRouter(newUserStore())

	body, _ := json.Marshal(createUserRequest{Name: "Antônio Prado", Email: "listagem@antonioeprado.dev"})
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var res []userResponse
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 user, got %d", len(res))
	}
}

func TestGetUserById(t *testing.T) {
	router := newRouter(newUserStore())

	body, _ := json.Marshal(createUserRequest{Name: "Antônio Prado", Email: "busca@antonioeprado.dev"})
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)))

	var created userResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/users/"+created.ID, nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestGetUserById_NotFound(t *testing.T) {
	router := newRouter(newUserStore())

	req := httptest.NewRequest(http.MethodGet, "/users/unknown-id", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdateUser(t *testing.T) {
	router := newRouter(newUserStore())

	createBody, _ := json.Marshal(createUserRequest{Name: "Antônio Prado", Email: "atualiza@antonioeprado.dev"})
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(createBody)))

	var created userResponse
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	updateBody, _ := json.Marshal(updateUserRequest{Name: "Antônio Elias Prado", Email: "atualiza@antonioeprado.dev"})
	req := httptest.NewRequest(http.MethodPut, "/users/"+created.ID, bytes.NewReader(updateBody))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var updated userResponse
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if updated.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	router := newRouter(newUserStore())

	body, _ := json.Marshal(updateUserRequest{Name: "Antônio Elias Prado", Email: "atualiza@antonioeprado.dev"})
	req := httptest.NewRequest(http.MethodPut, "/users/unknown-id", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
