package eventsourcing

import (
	uuid "github.com/satori/go.uuid"
)

type AggregateInterface interface {
	GetEntityType() string
	GetEntityUuid() uuid.UUID
	SetVersion(version int)
	GetVersion() int
	GetEvents() *EventStream
	ProcessEvent(event EventInterface)
	GetSnapshotStrategy() int

	ToDataStorage() *DataStorage
	FromDataStorage(dataStorage DataStorage)

	CreateEventFromDataStorage(eventType string, storage DataStorage) (EventInterface, error)
}
