package events_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aantonioprado/go-architecture/event-driven/internal/shared/events"
)

type testEvent struct {
	kind string
}

func (e testEvent) Type() string { return e.kind }

func TestBus_PublishCallsSubscriber(t *testing.T) {
	bus := events.NewBus()

	var got events.Event
	var mu sync.Mutex

	bus.Subscribe("thing.happened", func(e events.Event) {
		mu.Lock()
		defer mu.Unlock()
		got = e
	})

	bus.Publish(testEvent{kind: "thing.happened"})
	bus.Wait()

	mu.Lock()
	defer mu.Unlock()

	if got == nil {
		t.Fatal("expected handler to be called")
	}

	if got.Type() != "thing.happened" {
		t.Errorf("expected type %q, got %q", "thing.happened", got.Type())
	}
}

func TestBus_PublishCallsAllSubscribersForType(t *testing.T) {
	bus := events.NewBus()

	var calls int64

	bus.Subscribe("thing.happened", func(e events.Event) { atomic.AddInt64(&calls, 1) })
	bus.Subscribe("thing.happened", func(e events.Event) { atomic.AddInt64(&calls, 1) })

	bus.Publish(testEvent{kind: "thing.happened"})
	bus.Wait()

	if got := atomic.LoadInt64(&calls); got != 2 {
		t.Fatalf("expected 2 calls, got %d", got)
	}
}

func TestBus_PublishDoesNotCallSubscriberOfOtherType(t *testing.T) {
	bus := events.NewBus()

	var calls int64

	bus.Subscribe("other.type", func(e events.Event) { atomic.AddInt64(&calls, 1) })

	bus.Publish(testEvent{kind: "thing.happened"})
	bus.Wait()

	if got := atomic.LoadInt64(&calls); got != 0 {
		t.Fatalf("expected 0 calls, got %d", got)
	}
}

func TestBus_HandlerPanicDoesNotAffectOthers(t *testing.T) {
	bus := events.NewBus()

	var calls int64

	bus.Subscribe("thing.happened", func(e events.Event) {
		panic("boom")
	})
	bus.Subscribe("thing.happened", func(e events.Event) { atomic.AddInt64(&calls, 1) })

	bus.Publish(testEvent{kind: "thing.happened"})
	bus.Wait()

	if got := atomic.LoadInt64(&calls); got != 1 {
		t.Fatalf("expected the non-panicking handler to still run, got %d calls", got)
	}
}

func TestBus_PublishDoesNotBlock(t *testing.T) {
	bus := events.NewBus()

	release := make(chan struct{})

	bus.Subscribe("thing.happened", func(e events.Event) {
		<-release
	})

	done := make(chan struct{})
	go func() {
		bus.Publish(testEvent{kind: "thing.happened"})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Publish blocked waiting for the handler to finish")
	}

	close(release)
	bus.Wait()
}
