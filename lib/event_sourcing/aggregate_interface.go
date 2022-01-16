package eventsourcing

import (
	"distributes_system/lib/datastorage"
	uuid "github.com/satori/go.uuid"
)

type AggregateInterface interface {
	GetEntityType() string
	GetEntityUuid() uuid.UUID
	SetVersion(version int)
	GetVersion() int
	GetEvents() *EventStream
	ProcessEvent(event EventInterface)
	CreateEventFromDataStorage(eventType string, storage datastorage.DataStorage) (EventInterface, error)
}
