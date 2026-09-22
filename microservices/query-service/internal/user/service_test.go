package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/user"
)

func TestReadService_Replicate(t *testing.T) {
	svc := user.NewReadService(user.NewInMemoryUserRepository())

	err := svc.Replicate(user.User{
		ID:        "1",
		Name:      "Antônio Prado",
		Email:     "antonio@antonioeprado.dev",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := svc.GetUser("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email != "antonio@antonioeprado.dev" {
		t.Errorf("expected email %q, got %q", "antonio@antonioeprado.dev", found.Email)
	}
}

func TestReadService_ListUsers(t *testing.T) {
	svc := user.NewReadService(user.NewInMemoryUserRepository())

	if err := svc.Replicate(user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	users, err := svc.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestReadService_GetUser_NotFound(t *testing.T) {
	svc := user.NewReadService(user.NewInMemoryUserRepository())

	if _, err := svc.GetUser("unknown-id"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestReadService_RemoveReplica(t *testing.T) {
	svc := user.NewReadService(user.NewInMemoryUserRepository())

	if err := svc.Replicate(user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := svc.RemoveReplica("1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.GetUser("1"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestReadService_RemoveReplica_NotFound(t *testing.T) {
	svc := user.NewReadService(user.NewInMemoryUserRepository())

	if err := svc.RemoveReplica("unknown-id"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
