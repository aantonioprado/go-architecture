package usecases

import (
	"errors"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/entities"
)

var (
	ErrEmailTaken   = errors.New("email already in use")
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository is the persistence port. It is defined here, in the use
// case ring, and implemented by the outer gateway.
type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
}

// UserOutputPort is the presentation port. It is defined here, in the use
// case ring, and implemented by the outer presenter.
type UserOutputPort interface {
	PresentUserCreated(output UserOutput)
	PresentError(err error)
}

// UserInputPort is the boundary the outer controller depends on. It is
// implemented by the interactor.
type UserInputPort interface {
	CreateUser(input CreateUserInput, output UserOutputPort)
}
