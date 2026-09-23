package events

type Event interface {
	Type() string
}

type Handler func(Event)
