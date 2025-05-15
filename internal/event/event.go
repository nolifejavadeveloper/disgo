package event

import "reflect"

type Event interface {
	Type() string
}

var registry map[string]reflect.Type = make(map[string]reflect.Type)

func register(e Event) {
	registry[e.Type()] = reflect.TypeOf(e).Elem()
}