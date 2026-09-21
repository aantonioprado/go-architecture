package memory

import (
	"sort"
	"sync"

	"github.com/aantonioprado/go-architecture/ddd/internal/domain/user"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*user.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*user.User),
	}
}

func (r *UserRepository) Save(u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[u.ID()] = u

	return nil
}

func (r *UserRepository) FindAll() ([]*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt().Before(users[j].CreatedAt())
	})

	return users, nil
}

func (r *UserRepository) FindByID(id string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, user.ErrUserNotFound
	}

	return u, nil
}

func (r *UserRepository) FindByEmail(email user.Email) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email() == email {
			return u, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (r *UserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return user.ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}
