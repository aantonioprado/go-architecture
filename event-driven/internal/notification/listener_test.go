package notification_test

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aantonioprado/go-architecture/event-driven/internal/notification"
	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/events"
	"github.com/aantonioprado/go-architecture/event-driven/internal/user"
)

func TestListener_HandlesUserCreated(t *testing.T) {
	bus := events.NewBus()
	notification.NewListener(bus)

	var buf bytes.Buffer
	original := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(original)

	bus.Publish(user.UserCreated{
		UserID:    "1",
		Name:      "Antônio Prado",
		Email:     "antonio@antonioeprado.dev",
		CreatedAt: time.Now(),
	})
	bus.Wait()

	output := buf.String()

	if !strings.Contains(output, "[notification]") {
		t.Fatalf("expected notification log line, got: %q", output)
	}

	if !strings.Contains(output, "antonio@antonioeprado.dev") {
		t.Fatalf("expected log to mention the created user's email, got: %q", output)
	}
}

func TestListener_IgnoresOtherEventTypes(t *testing.T) {
	bus := events.NewBus()
	notification.NewListener(bus)

	var buf bytes.Buffer
	original := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(original)

	bus.Subscribe("other.type", func(e events.Event) {})
	bus.Publish(otherEvent{})
	bus.Wait()

	if strings.Contains(buf.String(), "[notification]") {
		t.Fatalf("expected no notification log for an unrelated event type, got: %q", buf.String())
	}
}

type otherEvent struct{}

func (otherEvent) Type() string { return "other.type" }
