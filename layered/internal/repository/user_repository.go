package repository

import (
	"errors"
	"sort"
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

func (r *UserRepository) FindById(id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepository) FindAll() ([]*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*model.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.Before(users[j].CreatedAt)
	})

	return users, nil
}

func (r *UserRepository) Update(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return ErrUserNotFound
	}

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
