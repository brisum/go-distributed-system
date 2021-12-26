package eventsourcing

import (
	uuid "github.com/satori/go.uuid"
)

type AggregateInterface interface {
	GetEntityType() string
	GetEntityUuid() uuid.UUID
	GetEvents() *EventStream
}
