package service_test

import (
	"errors"
	"testing"

	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/domain"
	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/service"
)

type fakeUserRepository struct {
	users map[string]*domain.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{users: make(map[string]*domain.User)}
}

func (r *fakeUserRepository) Create(user *domain.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *fakeUserRepository) FindAll() ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *fakeUserRepository) FindById(id string) (*domain.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *fakeUserRepository) Update(user *domain.User) error {
	if _, ok := r.users[user.ID]; !ok {
		return domain.ErrUserNotFound
	}

	r.users[user.ID] = user

	return nil
}

func (r *fakeUserRepository) Delete(id string) error {
	if _, ok := r.users[id]; !ok {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *fakeUserRepository) FindByEmail(email string) (*domain.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

func TestUserService_CreateUser(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepository())

	user, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Email != "antonio@antonioeprado.dev" {
		t.Errorf("expected email %q, got %q", "antonio@antonioeprado.dev", user.Email)
	}
}

func TestUserService_CreateUser_MissingName(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepository())

	if _, err := svc.CreateUser("", "antonio@antonioeprado.dev"); !errors.Is(err, domain.ErrNameRequired) {
		t.Fatalf("expected ErrNameRequired, got %v", err)
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepository())

	if _, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.CreateUser("Outro", "antonio@antonioeprado.dev"); !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserService_ListUsers(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepository())

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
	svc := service.NewUserService(newFakeUserRepository())

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
	svc := service.NewUserService(newFakeUserRepository())

	if _, err := svc.GetUser("unknown-id"); !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepository())

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
	svc := service.NewUserService(newFakeUserRepository())

	if _, err := svc.UpdateUser("unknown-id", "Antônio Prado", "antonio@antonioeprado.dev"); !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_UpdateUser_DuplicateEmail(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepository())

	if _, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := svc.CreateUser("Outro", "outro@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.UpdateUser(second.ID, "Outro", "antonio@antonioeprado.dev"); !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepository())

	created, err := svc.CreateUser("Antônio Prado", "antonio@antonioeprado.dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := svc.DeleteUser(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := svc.GetUser(created.ID); !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	svc := service.NewUserService(newFakeUserRepository())

	if err := svc.DeleteUser("unknown-id"); !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
