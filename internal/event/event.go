package event

import (
	"errors"
	"reflect"
)

type Event interface {
	Type() string
}

var registry map[string]reflect.Type = make(map[string]reflect.Type)

func register(e Event) {
	registry[e.Type()] = reflect.TypeOf(e).Elem()
}

func FindType(t string) (Event, error) {
	typ, ok := registry[t]
	if !ok {
		return nil, errors.New("event type not found in registry")
	}

	return reflect.New(typ).Interface().(Event), nil
}