package model

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"
)

var (
	ErrNameRequired  = errors.New("name is required")
	ErrEmailRequired = errors.New("email is required")
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(name, email string) (*User, error) {
	if err := validate(name, email); err != nil {
		return nil, err
	}

	return &User{
		ID:        newID(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}, nil
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
