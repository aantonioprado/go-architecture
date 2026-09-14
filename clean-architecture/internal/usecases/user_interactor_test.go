package usecases_test

import (
	"testing"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/entities"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

// fakeUserRepository and fakePresenter let the use case be tested without a
// real gateway or HTTP presenter, demonstrating the payoff of depending on
// interfaces rather than concrete types.
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
	err     error
}

func (p *fakePresenter) PresentUserCreated(output usecases.UserOutput) {
	p.created = &output
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
