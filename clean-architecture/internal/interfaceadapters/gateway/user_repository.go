package gateway

import (
	"sync"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/entities"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

// InMemoryUserRepository implements usecases.UserRepository. It is an outer
// ring type: it depends inward on entities and on the usecases port it
// satisfies, never the other way around.
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*entities.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*entities.User),
	}
}

func (r *InMemoryUserRepository) Create(user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user

	return nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, usecases.ErrUserNotFound
}
