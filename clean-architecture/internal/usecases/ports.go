package usecases

import (
	"errors"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/entities"
)

var (
	ErrEmailTaken   = errors.New("email already in use")
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Create(user *entities.User) error
	FindAll() ([]*entities.User, error)
	FindById(id string) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id string) error
}

type UserOutputPort interface {
	PresentUserCreated(output UserOutput)
	PresentUser(output UserOutput)
	PresentUserList(output ListUsersOutput)
	PresentUserDeleted()
	PresentError(err error)
}

type UserInputPort interface {
	CreateUser(input CreateUserInput, output UserOutputPort)
	ListUsers(output UserOutputPort)
	GetUserById(input GetUserInput, output UserOutputPort)
	UpdateUser(input UpdateUserInput, output UserOutputPort)
	DeleteUser(input DeleteUserInput, output UserOutputPort)
}
