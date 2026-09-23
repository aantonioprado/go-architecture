package audit

import (
	"log"

	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/events"
	"github.com/aantonioprado/go-architecture/event-driven/internal/user"
)

type Listener struct{}

func NewListener(bus *events.Bus) *Listener {
	l := &Listener{}

	bus.Subscribe(user.EventUserCreated, l.handle)

	return l
}

func (l *Listener) handle(e events.Event) {
	log.Printf("[audit] %s: %+v", e.Type(), e)
}
