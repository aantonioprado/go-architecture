package ports

import "github.com/aantonioprado/go-architecture/hexagonal/internal/core/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindAll() ([]*domain.User, error)
	FindById(id string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
}
