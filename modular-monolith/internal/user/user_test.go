package user_test

import (
	"errors"
	"testing"

	"github.com/aantonioprado/go-architecture/modular-monolith/internal/user"
)

func TestNewUser(t *testing.T) {
	u, err := user.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.ID == "" {
		t.Error("expected generated ID, got empty string")
	}

	if u.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if u.Name != "Antônio Prado" {
		t.Errorf("expected name %q, got %q", "Antônio Prado", u.Name)
	}
}

func TestNewUser_MissingName(t *testing.T) {
	if _, err := user.NewUser("", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestNewUser_MissingEmail(t *testing.T) {
	if _, err := user.NewUser("Antônio Prado", ""); !errors.Is(err, user.ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestUser_Update(t *testing.T) {
	u, err := user.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := u.Update("Antônio Elias Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.ID != u.ID {
		t.Errorf("expected ID to be preserved, got %q vs %q", updated.ID, u.ID)
	}

	if !updated.CreatedAt.Equal(u.CreatedAt) {
		t.Errorf("expected CreatedAt to be preserved, got %v vs %v", updated.CreatedAt, u.CreatedAt)
	}

	if updated.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}
}

func TestUser_Update_MissingName(t *testing.T) {
	u, err := user.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := u.Update("", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestUser_Update_MissingEmail(t *testing.T) {
	u, err := user.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := u.Update("Antônio Prado", ""); !errors.Is(err, user.ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}
