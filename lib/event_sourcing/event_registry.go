package eventsourcing

import "reflect"

type EventRegistry struct {
	eventRegistry map[string]reflect.Type
}

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		eventRegistry: make(map[string]reflect.Type, 0),
	}
}

func (registry *EventRegistry) Register(eventType string, reflectType reflect.Type) {
	registry.eventRegistry[eventType] = reflectType
}

func (registry *EventRegistry) GetReflectType(eventType string) reflect.Type {
	return registry.eventRegistry[eventType]
}
