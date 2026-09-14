package gateway_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/entities"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/interfaceadapters/gateway"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

func TestInMemoryUserRepository_Create(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	user := &entities.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}

	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInMemoryUserRepository_FindByEmail(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	user := &entities.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByEmail(user.Email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("expected id %q, got %q", user.ID, found.ID)
	}
}

func TestInMemoryUserRepository_FindByEmail_NotFound(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	if _, err := repo.FindByEmail("missing@antonioeprado.dev"); !errors.Is(err, usecases.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
