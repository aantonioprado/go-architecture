package user

import "time"

type Event interface {
	OccurredAt() time.Time
}

type UserRegistered struct {
	UserID string
	Name   string
	Email  string
	At     time.Time
}

func (e UserRegistered) OccurredAt() time.Time {
	return e.At
}

type UserDetailsChanged struct {
	UserID string
	Name   string
	Email  string
	At     time.Time
}

func (e UserDetailsChanged) OccurredAt() time.Time {
	return e.At
}

type UserDeleted struct {
	UserID string
	At     time.Time
}

func (e UserDeleted) OccurredAt() time.Time {
	return e.At
}
