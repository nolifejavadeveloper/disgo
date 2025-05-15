package pkg

import "github.com/nolifejavadeveloper/disgo/internal/event"


type EventBus struct {
	bus *event.Bus
}

func Subscribe[T any](b *EventBus, handler func(T) error) {
	event.Subscribe(b.bus, handler)
}

