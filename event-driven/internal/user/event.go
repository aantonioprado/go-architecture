package user

import (
	"time"

	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/events"
)

const EventUserCreated = "user.created"

type UserCreated struct {
	UserID    string
	Name      string
	Email     string
	CreatedAt time.Time
}

func (e UserCreated) Type() string {
	return EventUserCreated
}

type EventPublisher interface {
	Publish(event events.Event)
}
