package ports

import "github.com/aantonioprado/go-architecture/hexagonal/internal/core/domain"

// UserRepository is the secondary (driven) port: what the core needs from persistence.
type UserRepository interface {
	Create(user *domain.User) error
	FindAll() ([]*domain.User, error)
	FindById(id string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
}
