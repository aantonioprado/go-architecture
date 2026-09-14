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

func TestInMemoryUserRepository_FindAll(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	first := &entities.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	second := &entities.User{ID: "2", Name: "Another User", Email: "another@antonioeprado.dev", CreatedAt: time.Now()}

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

func TestInMemoryUserRepository_FindById(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	user := &entities.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
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

func TestInMemoryUserRepository_FindById_NotFound(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	if _, err := repo.FindById("unknown-id"); !errors.Is(err, usecases.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestInMemoryUserRepository_Update(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	user := &entities.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated := &entities.User{ID: user.ID, Name: "Antônio Elias Prado", Email: user.Email, CreatedAt: user.CreatedAt}
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

func TestInMemoryUserRepository_Update_NotFound(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	user := &entities.User{ID: "unknown-id", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}

	if err := repo.Update(user); !errors.Is(err, usecases.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestInMemoryUserRepository_Delete(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	user := &entities.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Delete(user.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := repo.FindById(user.ID); !errors.Is(err, usecases.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestInMemoryUserRepository_Delete_NotFound(t *testing.T) {
	repo := gateway.NewInMemoryUserRepository()

	if err := repo.Delete("unknown-id"); !errors.Is(err, usecases.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
