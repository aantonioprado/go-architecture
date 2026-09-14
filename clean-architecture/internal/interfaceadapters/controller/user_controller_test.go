package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/controller"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/dto"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/gateway"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

func newUserController() *controller.UserController {
	repo := gateway.NewInMemoryUserRepository()
	interactor := usecases.NewUserInteractor(repo)

	return controller.NewUserController(interactor)
}

func TestUserController_CreateUser(t *testing.T) {
	ctrl := newUserController()

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
}

func TestUserController_CreateUser_MissingName(t *testing.T) {
	ctrl := newUserController()

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
	ctrl := newUserController()

	body, _ := json.Marshal(dto.CreateUserRequest{
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
