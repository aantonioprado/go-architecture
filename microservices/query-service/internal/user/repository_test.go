package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/aantonioprado/go-architecture/microservices/query-service/internal/user"
)

func TestInMemoryUserRepository_Upsert(t *testing.T) {
	repo := user.NewInMemoryUserRepository()

	u := &user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}

	if err := repo.Upsert(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email != u.Email {
		t.Errorf("expected email %q, got %q", u.Email, found.Email)
	}
}

func TestInMemoryUserRepository_Upsert_Overwrites(t *testing.T) {
	repo := user.NewInMemoryUserRepository()

	u := &user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	if err := repo.Upsert(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated := &user.User{ID: "1", Name: "Antônio Elias Prado", Email: u.Email, CreatedAt: u.CreatedAt}
	if err := repo.Upsert(updated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByID("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", found.Name)
	}
}

func TestInMemoryUserRepository_FindAll(t *testing.T) {
	repo := user.NewInMemoryUserRepository()

	first := &user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	second := &user.User{ID: "2", Name: "Another User", Email: "another@antonioeprado.dev", CreatedAt: time.Now()}

	if err := repo.Upsert(first); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Upsert(second); err != nil {
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

func TestInMemoryUserRepository_FindByID_NotFound(t *testing.T) {
	repo := user.NewInMemoryUserRepository()

	if _, err := repo.FindByID("unknown-id"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestInMemoryUserRepository_Delete(t *testing.T) {
	repo := user.NewInMemoryUserRepository()

	u := &user.User{ID: "1", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev", CreatedAt: time.Now()}
	if err := repo.Upsert(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Delete(u.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := repo.FindByID(u.ID); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestInMemoryUserRepository_Delete_NotFound(t *testing.T) {
	repo := user.NewInMemoryUserRepository()

	if err := repo.Delete("unknown-id"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
