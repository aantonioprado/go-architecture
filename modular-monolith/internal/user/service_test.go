package user_test

import (
	"errors"
	"testing"

	"github.com/aantonioprado/go-architecture/modular-monolith/internal/user"
)

type fakeRepository struct {
	users map[string]*user.User
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{users: make(map[string]*user.User)}
}

func (r *fakeRepository) Create(u *user.User) error {
	r.users[u.ID] = u
	return nil
}

func (r *fakeRepository) FindAll() ([]*user.User, error) {
	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}

	return users, nil
}

func (r *fakeRepository) FindByID(id string) (*user.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}

	return u, nil
}

func (r *fakeRepository) Update(u *user.User) error {
	if _, ok := r.users[u.ID]; !ok {
		return user.ErrUserNotFound
	}

	r.users[u.ID] = u

	return nil
}

func (r *fakeRepository) Delete(id string) error {
	if _, ok := r.users[id]; !ok {
		return user.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *fakeRepository) FindByEmail(email string) (*user.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func TestUserService_CreateUser(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	u, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if u.Email != "antonio@antonioeprado.dev" {
		t.Errorf("expected email %q, got %q", "antonio@antonioeprado.dev", u.Email)
	}
}

func TestUserService_CreateUser_MissingName(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	if _, err := svc.CreateUser("", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	if _, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.CreateUser("Outro", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserService_ListUsers(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	if _, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	users, err := svc.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestUserService_GetUser(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	created, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found, err := svc.GetUser(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.Email != created.Email {
		t.Errorf("expected email %q, got %q", created.Email, found.Email)
	}
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	if _, err := svc.GetUser("unknown-id"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	created, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := svc.UpdateUser(created.ID, "Antônio Elias Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", updated.Name)
	}

	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("expected CreatedAt to be preserved, got %v vs %v", updated.CreatedAt, created.CreatedAt)
	}
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	if _, err := svc.UpdateUser("unknown-id", "Antônio Prado", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_UpdateUser_DuplicateEmail(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	if _, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := svc.CreateUser("Outro", "outro@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.UpdateUser(second.ID, "Outro", "antonio@antonioeprado.dev"); !errors.Is(err, user.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	created, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := svc.DeleteUser(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.GetUser(created.ID); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	svc := user.NewUserService(newFakeRepository())

	if err := svc.DeleteUser("unknown-id"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
