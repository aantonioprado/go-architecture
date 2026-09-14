package repository

import (
	"errors"
	"sync"

	"github.com/aantonioprado/go-architecture/layered/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*model.User),
	}
}

func (r *UserRepository) Create(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user

	return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}
