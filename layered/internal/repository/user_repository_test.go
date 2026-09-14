package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aantonioprado/go-architecture/layered/internal/model"
	"github.com/aantonioprado/go-architecture/layered/internal/repository"
)

func TestUserRepository_Create(t *testing.T) {
	repo := repository.NewUserRepository()

	user := &model.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}

	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	repo := repository.NewUserRepository()

	user := &model.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
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

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	repo := repository.NewUserRepository()

	if _, err := repo.FindByEmail("missing@antonioeprado.dev"); !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
