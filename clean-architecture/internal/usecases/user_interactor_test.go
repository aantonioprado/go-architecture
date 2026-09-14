package usecases_test

import (
	"testing"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/entities"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

type fakeUserRepository struct {
	users map[string]*entities.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{users: make(map[string]*entities.User)}
}

func (r *fakeUserRepository) Create(user *entities.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *fakeUserRepository) FindAll() ([]*entities.User, error) {
	users := make([]*entities.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *fakeUserRepository) FindById(id string) (*entities.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, usecases.ErrUserNotFound
	}

	return user, nil
}

func (r *fakeUserRepository) Update(user *entities.User) error {
	if _, ok := r.users[user.ID]; !ok {
		return usecases.ErrUserNotFound
	}

	r.users[user.ID] = user

	return nil
}

func (r *fakeUserRepository) Delete(id string) error {
	if _, ok := r.users[id]; !ok {
		return usecases.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *fakeUserRepository) FindByEmail(email string) (*entities.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, usecases.ErrUserNotFound
}

type fakePresenter struct {
	created *usecases.UserOutput
	user    *usecases.UserOutput
	list    *usecases.ListUsersOutput
	deleted bool
	err     error
}

func (p *fakePresenter) PresentUserCreated(output usecases.UserOutput) {
	p.created = &output
}

func (p *fakePresenter) PresentUser(output usecases.UserOutput) {
	p.user = &output
}

func (p *fakePresenter) PresentUserList(output usecases.ListUsersOutput) {
	p.list = &output
}

func (p *fakePresenter) PresentUserDeleted() {
	p.deleted = true
}

func (p *fakePresenter) PresentError(err error) {
	p.err = err
}

func TestUserInteractor_CreateUser(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())
	out := &fakePresenter{}

	interactor.CreateUser(usecases.CreateUserInput{
		Name:  "Antônio Prado",
		Email: "antonio@antonioeprado.dev",
	}, out)

	if out.err != nil {
		t.Fatalf("unexpected error: %v", out.err)
	}

	if out.created == nil {
		t.Fatal("expected PresentUserCreated to be called")
	}

	if out.created.Email != "antonio@antonioeprado.dev" {
		t.Errorf("expected email %q, got %q", "antonio@antonioeprado.dev", out.created.Email)
	}
}

func TestUserInteractor_CreateUser_MissingName(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())
	out := &fakePresenter{}

	interactor.CreateUser(usecases.CreateUserInput{
		Email: "antonio@antonioeprado.dev",
	}, out)

	if out.err == nil {
		t.Fatal("expected PresentError to be called")
	}
}

func TestUserInteractor_CreateUser_DuplicateEmail(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())

	first := &fakePresenter{}
	interactor.CreateUser(usecases.CreateUserInput{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"}, first)

	second := &fakePresenter{}
	interactor.CreateUser(usecases.CreateUserInput{Name: "Outro", Email: "antonio@antonioeprado.dev"}, second)

	if second.err == nil {
		t.Fatal("expected PresentError to be called for duplicate email")
	}
}

func TestUserInteractor_ListUsers(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())

	interactor.CreateUser(usecases.CreateUserInput{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"}, &fakePresenter{})

	out := &fakePresenter{}
	interactor.ListUsers(out)

	if out.list == nil {
		t.Fatal("expected PresentUserList to be called")
	}

	if len(out.list.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(out.list.Users))
	}
}

func TestUserInteractor_GetUserById(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())

	created := &fakePresenter{}
	interactor.CreateUser(usecases.CreateUserInput{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"}, created)

	out := &fakePresenter{}
	interactor.GetUserById(usecases.GetUserInput{ID: created.created.ID}, out)

	if out.user == nil {
		t.Fatal("expected PresentUser to be called")
	}
}

func TestUserInteractor_GetUserById_NotFound(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())

	out := &fakePresenter{}
	interactor.GetUserById(usecases.GetUserInput{ID: "unknown-id"}, out)

	if out.err == nil {
		t.Fatal("expected PresentError to be called")
	}
}

func TestUserInteractor_UpdateUser(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())

	created := &fakePresenter{}
	interactor.CreateUser(usecases.CreateUserInput{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"}, created)

	out := &fakePresenter{}
	interactor.UpdateUser(usecases.UpdateUserInput{
		ID:    created.created.ID,
		Name:  "Antônio Elias Prado",
		Email: "antonio@antonioeprado.dev",
	}, out)

	if out.user == nil {
		t.Fatal("expected PresentUser to be called")
	}

	if out.user.Name != "Antônio Elias Prado" {
		t.Errorf("expected updated name, got %q", out.user.Name)
	}

	if !out.user.CreatedAt.Equal(created.created.CreatedAt) {
		t.Errorf("expected CreatedAt to be preserved, got %v vs %v", out.user.CreatedAt, created.created.CreatedAt)
	}
}

func TestUserInteractor_UpdateUser_NotFound(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())

	out := &fakePresenter{}
	interactor.UpdateUser(usecases.UpdateUserInput{ID: "unknown-id", Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"}, out)

	if out.err == nil {
		t.Fatal("expected PresentError to be called")
	}
}

func TestUserInteractor_UpdateUser_DuplicateEmail(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())

	first := &fakePresenter{}
	interactor.CreateUser(usecases.CreateUserInput{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"}, first)

	second := &fakePresenter{}
	interactor.CreateUser(usecases.CreateUserInput{Name: "Outro", Email: "outro@antonioeprado.dev"}, second)

	out := &fakePresenter{}
	interactor.UpdateUser(usecases.UpdateUserInput{ID: second.created.ID, Name: "Outro", Email: "antonio@antonioeprado.dev"}, out)

	if out.err == nil {
		t.Fatal("expected PresentError to be called for duplicate email")
	}
}

func TestUserInteractor_DeleteUser(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())

	created := &fakePresenter{}
	interactor.CreateUser(usecases.CreateUserInput{Name: "Antônio Prado", Email: "antonio@antonioeprado.dev"}, created)

	out := &fakePresenter{}
	interactor.DeleteUser(usecases.DeleteUserInput{ID: created.created.ID}, out)

	if !out.deleted {
		t.Fatal("expected PresentUserDeleted to be called")
	}
}

func TestUserInteractor_DeleteUser_NotFound(t *testing.T) {
	interactor := usecases.NewUserInteractor(newFakeUserRepository())

	out := &fakePresenter{}
	interactor.DeleteUser(usecases.DeleteUserInput{ID: "unknown-id"}, out)

	if out.err == nil {
		t.Fatal("expected PresentError to be called")
	}
}
