package event

import (
	"reflect"
)

type EventHandler func(any) error

type Bus struct {
	handlers map[reflect.Type][]EventHandler
}

func NewBus() *Bus {
	return &Bus{
		handlers: make(map[reflect.Type][]EventHandler),
	}
}

func Subscribe[T any](b *Bus, handler func(T) error) {
	var t T
	typ := reflect.TypeOf(t)
	b.handlers[typ] = append(b.handlers[typ], func(e any) error {
		return handler(e.(T))
	})
}

func Fire(b *Bus, event any) {
	for _, handler := range b.handlers[reflect.TypeOf(event)] {
		handler(event)
	}
}
