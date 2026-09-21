package user

import (
	"crypto/rand"
	"fmt"
	"time"
)

type User struct {
	id        string
	name      string
	email     Email
	createdAt time.Time
	events    []Event
}

func Register(name string, email Email) (*User, error) {
	if name == "" {
		return nil, ErrNameRequired
	}

	u := &User{
		id:        newID(),
		name:      name,
		email:     email,
		createdAt: time.Now(),
	}

	u.record(UserRegistered{UserID: u.id, Name: u.name, Email: u.email.String(), At: u.createdAt})

	return u, nil
}

func (u *User) ChangeDetails(name string, email Email) error {
	if name == "" {
		return ErrNameRequired
	}

	u.name = name
	u.email = email

	u.record(UserDetailsChanged{UserID: u.id, Name: u.name, Email: u.email.String(), At: time.Now()})

	return nil
}

func (u *User) record(event Event) {
	u.events = append(u.events, event)
}

func (u *User) Events() []Event {
	return u.events
}

func (u *User) ClearEvents() {
	u.events = nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
