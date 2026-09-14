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

func TestUserRepository_FindById(t *testing.T) {
	repo := repository.NewUserRepository()

	user := &model.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindById(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email != user.Email {
		t.Errorf("expected email %q, got %q", user.Email, found.Email)
	}
}

func TestUserRepository_FindById_NotFound(t *testing.T) {
	repo := repository.NewUserRepository()

	if _, err := repo.FindById("unknown-id"); !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	repo := repository.NewUserRepository()

	first := &model.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	second := &model.User{ID: "2", Name: "Another User", Email: "another@antonioeprado.dev", CreatedAt: time.Now()}

	if err := repo.Create(first); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Create(second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	users, err := repo.FindAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := repository.NewUserRepository()

	user := &model.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated := &model.User{ID: user.ID, Name: "Antônio Elias Prado", Email: user.Email, CreatedAt: user.CreatedAt}
	if err := repo.Update(updated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindById(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", found.Name)
	}
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	repo := repository.NewUserRepository()

	user := &model.User{ID: "unknown-id", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}

	if err := repo.Update(user); !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
