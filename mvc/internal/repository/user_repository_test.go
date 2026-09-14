package repository_test

import (
	"errors"
	"testing"

	"github.com/aantonioprado/go-architecture/mvc/internal/model"
	"github.com/aantonioprado/go-architecture/mvc/internal/repository"
)

func TestUserRepository_Create(t *testing.T) {
	repo := repository.NewUserRepository()

	user, err := model.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error creating user: %v", err)
	}

	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	repo := repository.NewUserRepository()

	first, _ := model.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	second, _ := model.NewUser("Another User", "antonio@antonioeprado.dev")

	if err := repo.Create(first); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Create(second); !errors.Is(err, repository.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserRepository_FindById(t *testing.T) {
	repo := repository.NewUserRepository()

	user, _ := model.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
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

	first, _ := model.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	second, _ := model.NewUser("Another User", "another@antonioeprado.dev")

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

	user, _ := model.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := user.Update("Antônio Elias Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

	user, _ := model.NewUser("Antônio Prado", "antonio@antonioeprado.dev")

	if err := repo.Update(user); !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_Update_DuplicateEmail(t *testing.T) {
	repo := repository.NewUserRepository()

	first, _ := model.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	second, _ := model.NewUser("Another User", "another@antonioeprado.dev")

	if err := repo.Create(first); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Create(second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := second.Update(second.Name, first.Email)

	if err := repo.Update(updated); !errors.Is(err, repository.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo := repository.NewUserRepository()

	user, _ := model.NewUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err := repo.Create(user); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Delete(user.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := repo.FindById(user.ID); !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	repo := repository.NewUserRepository()

	if err := repo.Delete("unknown-id"); !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
