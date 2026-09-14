package gateway

import (
	"sort"
	"sync"

	"github.com/aantonioprado/go-architecture/clean-architecture/internal/entities"
	"github.com/aantonioprado/go-architecture/clean-architecture/internal/usecases"
)

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

func (r *InMemoryUserRepository) FindAll() ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entities.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.Before(users[j].CreatedAt)
	})

	return users, nil
}

func (r *InMemoryUserRepository) FindById(id string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, usecases.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) Update(user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return usecases.ErrUserNotFound
	}

	r.users[user.ID] = user

	return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return usecases.ErrUserNotFound
	}

	delete(r.users, id)

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
