package repository

import (
	"errors"
	"sync"

	"github.com/aantonioprado/go-architecture/mvc/internal/model"
)

var ErrEmailTaken = errors.New("email already in use")

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

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	return nil, nil
}

func (r *UserRepository) FindAll() ([]*model.User, error) {
	return nil, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return nil
}

func (r *UserRepository) Delete(id string) error {
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
