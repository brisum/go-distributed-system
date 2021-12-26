package eventsourcing

import (
	uuid "github.com/satori/go.uuid"
)

type EventStream struct {
	entityType string
	entityUuid uuid.UUID
	events     []EventInterface
}

func NewEventStream(entityType string, entityUuid uuid.UUID) *EventStream {
	stream := EventStream{
		entityType: entityType,
		entityUuid: entityUuid,
		events:     make([]EventInterface, 0),
	}

	return &stream
}

func (stream *EventStream) GetEntityType() string {
	return stream.entityType
}

func (stream *EventStream) GetEntityUuid() uuid.UUID {
	return stream.entityUuid
}

func (stream *EventStream) AppendEvent(event EventInterface) {
	stream.events = append(stream.events, event)
}

func (stream *EventStream) GetEvents() []EventInterface {
	events := stream.events

	stream.events = make([]EventInterface, 0)

	return events
}
