package memory

import (
	"sort"
	"sync"

	"github.com/aantonioprado/go-architecture/hexagonal/internal/core/domain"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user

	return nil
}

func (r *InMemoryUserRepository) FindAll() ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.Before(users[j].CreatedAt)
	})

	return users, nil
}

func (r *InMemoryUserRepository) FindById(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

func (r *InMemoryUserRepository) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return domain.ErrUserNotFound
	}

	r.users[user.ID] = user

	return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *InMemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}
