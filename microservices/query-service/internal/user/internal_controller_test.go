package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/user"
)

func TestInternalController_Replicate(t *testing.T) {
	repo := user.NewInMemoryUserRepository()
	svc := user.NewReadService(repo)
	ctrl := user.NewInternalController(svc)

	body, _ := json.Marshal(user.ReplicaRequest{
		ID:        "1",
		Name:      "Antônio Prado",
		Email:     "antonio@antonioeprado.dev",
		CreatedAt: time.Now().Format(time.RFC3339),
	})

	req := httptest.NewRequest(http.MethodPost, "/internal/users", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ctrl.Replicate(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	found, err := svc.GetUser("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email != "antonio@antonioeprado.dev" {
		t.Errorf("expected email %q, got %q", "antonio@antonioeprado.dev", found.Email)
	}
}

func TestInternalController_UpdateReplica(t *testing.T) {
	repo := user.NewInMemoryUserRepository()
	svc := user.NewReadService(repo)
	ctrl := user.NewInternalController(svc)

	if err := svc.Replicate(user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	body, _ := json.Marshal(user.ReplicaRequest{
		Name:      "Antônio Elias Prado",
		Email:     "antonio@antonioeprado.dev",
		CreatedAt: time.Now().Format(time.RFC3339),
	})

	req := newRequestWithID(http.MethodPut, "/internal/users/1", "1", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ctrl.UpdateReplica(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	found, err := svc.GetUser("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", found.Name)
	}
}

func TestInternalController_RemoveReplica(t *testing.T) {
	repo := user.NewInMemoryUserRepository()
	svc := user.NewReadService(repo)
	ctrl := user.NewInternalController(svc)

	if err := svc.Replicate(user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := newRequestWithID(http.MethodDelete, "/internal/users/1", "1", nil)
	rec := httptest.NewRecorder()

	ctrl.RemoveReplica(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if _, err := svc.GetUser("1"); err == nil {
		t.Fatal("expected user to have been removed")
	}
}
