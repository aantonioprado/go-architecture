package audit_test

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aantonioprado/go-architecture/event-driven/internal/audit"
	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/events"
	"github.com/aantonioprado/go-architecture/event-driven/internal/user"
)

func TestListener_HandlesUserCreated(t *testing.T) {
	bus := events.NewBus()
	audit.NewListener(bus)

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

	if !strings.Contains(output, "[audit]") {
		t.Fatalf("expected audit log line, got: %q", output)
	}

	if !strings.Contains(output, user.EventUserCreated) {
		t.Fatalf("expected log to mention the event type, got: %q", output)
	}
}
