package memory_test

import (
	"errors"
	"testing"

	"github.com/aantonioprado/go-architecture/ddd/internal/domain/user"
	"github.com/aantonioprado/go-architecture/ddd/internal/infrastructure/persistence/memory"
)

func mustUser(t *testing.T, name, rawEmail string) *user.User {
	t.Helper()

	email, err := user.NewEmail(rawEmail)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, err := user.Register(name, email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return u
}

func TestUserRepository_Save(t *testing.T) {
	repo := memory.NewUserRepository()

	u := mustUser(t, "Antônio Prado", "antonio@antonioeprado.dev")

	if err := repo.Save(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	repo := memory.NewUserRepository()

	u := mustUser(t, "Antônio Prado", "antonio@antonioeprado.dev")
	if err := repo.Save(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByEmail(u.Email())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.ID() != u.ID() {
		t.Errorf("expected id %q, got %q", u.ID(), found.ID())
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	repo := memory.NewUserRepository()

	email, err := user.NewEmail("missing@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := repo.FindByEmail(email); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_FindAll(t *testing.T) {
	repo := memory.NewUserRepository()

	first := mustUser(t, "Antônio Prado", "antonio@antonioeprado.dev")
	second := mustUser(t, "Another User", "another@antonioeprado.dev")

	if err := repo.Save(first); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Save(second); err != nil {
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

func TestUserRepository_FindByID(t *testing.T) {
	repo := memory.NewUserRepository()

	u := mustUser(t, "Antônio Prado", "antonio@antonioeprado.dev")
	if err := repo.Save(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByID(u.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email() != u.Email() {
		t.Errorf("expected email %q, got %q", u.Email(), found.Email())
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	repo := memory.NewUserRepository()

	if _, err := repo.FindByID("unknown-id"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo := memory.NewUserRepository()

	u := mustUser(t, "Antônio Prado", "antonio@antonioeprado.dev")
	if err := repo.Save(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Delete(u.ID()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := repo.FindByID(u.ID()); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	repo := memory.NewUserRepository()

	if err := repo.Delete("unknown-id"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
