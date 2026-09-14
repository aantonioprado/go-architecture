package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNameRequired  = errors.New("name is required")
	ErrEmailRequired = errors.New("email is required")
	ErrEmailTaken    = errors.New("email already in use")
	ErrUserNotFound  = errors.New("user not found")
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

type userStore struct {
	mu    sync.RWMutex
	users map[string]*User
}

func newUserStore() *userStore {
	return &userStore{
		users: make(map[string]*User),
	}
}

func (s *userStore) create(name, email string) (*User, error) {
	if err := validate(name, email); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.findByEmail(email) != nil {
		return nil, ErrEmailTaken
	}

	user := &User{
		ID:        newID(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	s.users[user.ID] = user

	return user, nil
}

func (s *userStore) findByEmail(email string) *User {
	for _, user := range s.users {
		if user.Email == email {
			return user
		}
	}

	return nil
}

func validate(name, email string) error {
	if name == "" {
		return ErrNameRequired
	}

	if email == "" {
		return ErrEmailRequired
	}

	return nil
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
