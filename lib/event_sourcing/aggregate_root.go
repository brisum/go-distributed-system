package eventsourcing

import (
	uuid "github.com/satori/go.uuid"
)

type AggregateRoot struct {
	entityType        string
	entityUuid        uuid.UUID
	entityVersion     int
	snapshotStrategy  int
	uncommittedEvents []EventInterface
}

func NewAggregateRoot(entityType string, entityUuid uuid.UUID, snapshotStrategy int) *AggregateRoot {
	return &AggregateRoot{
		entityType:        entityType,
		entityUuid:        entityUuid,
		entityVersion:     0,
		snapshotStrategy:  snapshotStrategy,
		uncommittedEvents: make([]EventInterface, 0),
	}
}

func (aggregate *AggregateRoot) GetEntityType() string {
	return aggregate.entityType
}

func (aggregate *AggregateRoot) GetEntityUuid() uuid.UUID {
	return aggregate.entityUuid
}

func (aggregate *AggregateRoot) SetVersion(version int) {
	aggregate.entityVersion = version
}

func (aggregate *AggregateRoot) GetVersion() int {
	return aggregate.entityVersion
}

func (aggregate *AggregateRoot) GetSnapshotStrategy() int {
	return aggregate.snapshotStrategy
}

func (aggregate *AggregateRoot) AppendEvent(event EventInterface) {
	aggregate.uncommittedEvents = append(aggregate.uncommittedEvents, event)
}

func (aggregate *AggregateRoot) GetUncommittedEvents() *EventStream {
	stream := NewEventStream(aggregate.entityType, aggregate.entityUuid)
	reversedEvents := make([]EventInterface, 0)

	for _, event := range aggregate.uncommittedEvents {
		reversedEvents = append(reversedEvents, event)
	}

	for _, event := range reversedEvents {
		stream.AppendEvent(event)
	}

	aggregate.uncommittedEvents = make([]EventInterface, 0)

	return stream
}
