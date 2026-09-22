package user

import "errors"

var (
	ErrNameRequired  = errors.New("name is required")
	ErrEmailRequired = errors.New("email is required")
	ErrEmailTaken    = errors.New("email already in use")
	ErrUserNotFound  = errors.New("user not found")
)
