package user_test

import (
	"errors"
	"testing"

	"github.com/aantonioprado/go-architecture/ddd/internal/domain/user"
)

func TestNewEmail(t *testing.T) {
	email, err := user.NewEmail("antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if email.String() != "antonio@antonioeprado.dev" {
		t.Errorf("expected %q, got %q", "antonio@antonioeprado.dev", email.String())
	}
}

func TestNewEmail_Empty(t *testing.T) {
	if _, err := user.NewEmail(""); !errors.Is(err, user.ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}
