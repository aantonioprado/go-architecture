package domain_test

import (
	"errors"
	"testing"

	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/domain"
)

func TestNewUser(t *testing.T) {
	user, err := domain.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID == "" {
		t.Error("expected generated ID, got empty string")
	}

	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if user.Name != "Antônio Prado" {
		t.Errorf("expected name %q, got %q", "Antônio Prado", user.Name)
	}
}

func TestNewUser_MissingName(t *testing.T) {
	if _, err := domain.NewUser("", "antonio@antonioeprado.dev"); !errors.Is(err, domain.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestNewUser_MissingEmail(t *testing.T) {
	if _, err := domain.NewUser("Antônio Prado", ""); !errors.Is(err, domain.ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestUser_Update(t *testing.T) {
	user, err := domain.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := user.Update("Antônio Elias Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.ID != user.ID {
		t.Errorf("expected ID to be preserved, got %q vs %q", updated.ID, user.ID)
	}

	if !updated.CreatedAt.Equal(user.CreatedAt) {
		t.Errorf("expected CreatedAt to be preserved, got %v vs %v", updated.CreatedAt, user.CreatedAt)
	}

	if updated.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}
}

func TestUser_Update_MissingName(t *testing.T) {
	user, err := domain.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := user.Update("", "antonio@antonioeprado.dev"); !errors.Is(err, domain.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestUser_Update_MissingEmail(t *testing.T) {
	user, err := domain.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := user.Update("Antônio Prado", ""); !errors.Is(err, domain.ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}
