package service_test

import (
	"errors"
	"testing"

	"github.com/aantonioprado/go-architecture/layered/internal/repository"
	"github.com/aantonioprado/go-architecture/layered/internal/service"
)

func TestUserService_Create(t *testing.T) {
	svc := service.NewUserService(repository.NewUserRepository())

	user, err := svc.Create("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID == "" {
		t.Error("expected generated ID, got empty string")
	}

	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestUserService_Create_MissingName(t *testing.T) {
	svc := service.NewUserService(repository.NewUserRepository())

	if _, err := svc.Create("", "antonio@antonioeprado.dev"); !errors.Is(err, service.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestUserService_Create_MissingEmail(t *testing.T) {
	svc := service.NewUserService(repository.NewUserRepository())

	if _, err := svc.Create("Antônio Prado", ""); !errors.Is(err, service.ErrEmailRequired) {
		t.Fatalf("expected ErrEmailRequired, got %v", err)
	}
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	svc := service.NewUserService(repository.NewUserRepository())

	if _, err := svc.Create("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.Create("Another User", "antonio@antonioeprado.dev"); !errors.Is(err, service.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserService_GetById(t *testing.T) {
	svc := service.NewUserService(repository.NewUserRepository())

	created, err := svc.Create("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := svc.GetById(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email != created.Email {
		t.Errorf("expected email %q, got %q", created.Email, found.Email)
	}
}

func TestUserService_GetById_NotFound(t *testing.T) {
	svc := service.NewUserService(repository.NewUserRepository())

	if _, err := svc.GetById("unknown-id"); !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_List(t *testing.T) {
	svc := service.NewUserService(repository.NewUserRepository())

	if _, err := svc.Create("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	users, err := svc.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}
