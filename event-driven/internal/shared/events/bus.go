package events

import (
	"log"
	"sync"
)

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	wg       sync.WaitGroup
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]Handler),
	}
}

func (b *Bus) Subscribe(eventType string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

func (b *Bus) Publish(event Event) {
	b.mu.RLock()
	handlers := append([]Handler(nil), b.handlers[event.Type()]...)
	b.mu.RUnlock()

	for _, handler := range handlers {
		b.wg.Add(1)

		go func(h Handler, e Event) {
			defer b.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[event-bus] handler for %q panicked: %v", e.Type(), r)
				}
			}()

			h(e)
		}(handler, event)
	}
}

func (b *Bus) Wait() {
	b.wg.Wait()
}
