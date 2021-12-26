package eventsourcing

import (
	"distributes_system/lib/datastorage"
	uuid "github.com/satori/go.uuid"
)

type AggregateInterface interface {
	GetEntityType() string
	GetEntityUuid() uuid.UUID
	GetEvents() *EventStream
	ProcessEvent(event EventInterface)
	ApplyEvent(event EventInterface)
	CreateEventFromDataStorage(eventType string, storage datastorage.DataStorage) (EventInterface, error)
}
