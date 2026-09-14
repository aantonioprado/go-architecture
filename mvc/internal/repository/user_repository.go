package repository

import (
	"errors"
	"sort"
	"sync"

	"github.com/aantonioprado/go-architecture/mvc/internal/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailTaken   = errors.New("email already in use")
)

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

	if r.emailTaken(user.Email, "") {
		return ErrEmailTaken
	}

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

	if r.emailTaken(user.Email, user.ID) {
		return ErrEmailTaken
	}

	r.users[user.ID] = user

	return nil
}

func (r *UserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return ErrUserNotFound
	}

	delete(r.users, id)

	return nil
}

func (r *UserRepository) emailTaken(email, excludeID string) bool {
	for _, u := range r.users {
		if u.Email == email && u.ID != excludeID {
			return true
		}
	}

	return false
}
