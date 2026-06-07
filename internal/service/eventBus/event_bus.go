package eventBus

import "log"

type Bus struct {
	events chan Event
}

func NewBus(capacity int) *Bus {
	return &Bus{
		events: make(chan Event, capacity),
	}
}

func (b *Bus) Publish(event Event) {
	select {
	case b.events <- event:
	default:
		log.Println("WARN: event dropped, event bus buffer is full")
	}
}

func (b *Bus) Subscribe() <-chan Event {
	return b.events
}

func (b *Bus) Close() {
	close(b.events)
}
