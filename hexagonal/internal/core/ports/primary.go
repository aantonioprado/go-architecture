package ports

import "github.com/aantonioprado/go-architecture/hexagonal/internal/core/domain"

type UserService interface {
	CreateUser(name, email string) (domain.User, error)
	ListUsers() ([]domain.User, error)
	GetUser(id string) (domain.User, error)
	UpdateUser(id, name, email string) (domain.User, error)
	DeleteUser(id string) error
}
