package user_test

import (
	"errors"
	"testing"

	"github.com/aantonioprado/go-architecture/ddd/internal/domain/user"
)

func mustEmail(t *testing.T, raw string) user.Email {
	t.Helper()

	email, err := user.NewEmail(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return email
}

func TestRegister(t *testing.T) {
	email := mustEmail(t, "antonio@antonioeprado.dev")

	u, err := user.Register("Antônio Prado", email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.ID() == "" {
		t.Error("expected generated ID, got empty string")
	}

	if u.CreatedAt().IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if u.Name() != "Antônio Prado" {
		t.Errorf("expected name %q, got %q", "Antônio Prado", u.Name())
	}
}

func TestRegister_MissingName(t *testing.T) {
	email := mustEmail(t, "antonio@antonioeprado.dev")

	if _, err := user.Register("", email); !errors.Is(err, user.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestUser_ChangeDetails(t *testing.T) {
	email := mustEmail(t, "antonio@antonioeprado.dev")

	u, err := user.Register("Antônio Prado", email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	id := u.ID()
	createdAt := u.CreatedAt()

	if err := u.ChangeDetails("Antônio Elias Prado", email); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.ID() != id {
		t.Errorf("expected ID to be preserved, got %q vs %q", u.ID(), id)
	}

	if !u.CreatedAt().Equal(createdAt) {
		t.Errorf("expected CreatedAt to be preserved, got %v vs %v", u.CreatedAt(), createdAt)
	}

	if u.Name() != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", u.Name())
	}
}

func TestUser_ChangeDetails_MissingName(t *testing.T) {
	email := mustEmail(t, "antonio@antonioeprado.dev")

	u, err := user.Register("Antônio Prado", email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := u.ChangeDetails("", email); !errors.Is(err, user.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}
