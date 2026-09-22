package memory_test

import (
	"errors"
	"sync"
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

func TestUserRepository_FindByID_ReturnsIsolatedCopy(t *testing.T) {
	repo := memory.NewUserRepository()

	u := mustUser(t, "Antônio Prado", "antonio@antonioeprado.dev")
	if err := repo.Save(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := repo.FindByID(u.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	email, err := user.NewEmail("mutated@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := found.ChangeDetails("Mutated", email); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stillStored, err := repo.FindByID(u.ID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stillStored.Name() == "Mutated" {
		t.Fatal("expected mutating a returned user not to affect what the repository has stored")
	}
}

func TestUserRepository_ConcurrentReadWriteSameID(t *testing.T) {
	repo := memory.NewUserRepository()

	email, err := user.NewEmail("race@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, err := user.Register("Race User", email)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.Save(u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			if found, err := repo.FindByID(u.ID()); err == nil {
				_ = found.Name()
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 200; i++ {
			found, err := repo.FindByID(u.ID())
			if err != nil {
				continue
			}

			if err := found.ChangeDetails("Mutated", email); err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if err := repo.Save(found); err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
		}
	}()

	wg.Wait()
}
