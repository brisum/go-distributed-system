package eventsourcing

import (
	"context"
	"distributes_system/lib/datastorage"
	"github.com/georgysavva/scany/pgxscan"
	pgx "github.com/jackc/pgx/v4"
	"reflect"
)

type EventStore struct {
	connection    *pgx.Conn
	eventRegistry *EventRegistry
}

func NewEventStore(connection *pgx.Conn) *EventStore {
	store := EventStore{
		connection:    connection,
		eventRegistry: NewEventRegistry(),
	}

	return &store
}

func (store *EventStore) Save(ctx *context.Context, aggregate AggregateInterface) error {
	stream := aggregate.GetEvents()

	reversedEvents := make([]EventInterface, 0)
	for _, event := range stream.GetEvents() {
		reversedEvents = append(reversedEvents, event)
	}

	tx, err := store.connection.Begin(*ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(*ctx)

	for _, event := range reversedEvents {
		_, err = tx.Exec(
			*ctx,
			"INSERT INTO event (event_type, entity_type, entity_uuid, event_data, promoter, triggering_event) "+
				"VALUES ($1, $2, $3, $4, $5, $6)",
			GetEventType(event),
			aggregate.GetEntityType(),
			aggregate.GetEntityUuid().String(),
			event.ToDataStorage().MarshalJSON(),
			"",
			"",
		)

		if err != nil {
			return err
		}
	}

	err = tx.Commit(*ctx)
	if err != nil {
		return err
	}

	return nil
}

func (store *EventStore) Load(ctx *context.Context, aggregate AggregateInterface) error {
	eventRows := make([]EventRow, 0)

	err := pgxscan.Select(
		*ctx,
		store.connection,
		&eventRows,
		`SELECT event_type, event_data FROM event WHERE entity_type = $1 and entity_uuid = $2`,
		aggregate.GetEntityType(),
		aggregate.GetEntityUuid().String(),
	)

	if err != nil {
		return err
	}

	for _, eventRow := range eventRows {
		storage := datastorage.NewEmptyDataStorage()
		storage.UnmarshalJSON(eventRow.EventData)

		event, _ := aggregate.CreateEventFromDataStorage(eventRow.EventType, *storage)

		methodName := "Apply" + eventRow.EventType + "Event"
		methodArguments := []reflect.Value{reflect.ValueOf(event)}
		aggregateValue := reflect.ValueOf(aggregate)
		aggregateValue.MethodByName(methodName).Call(methodArguments)
	}

	return nil
}

func (store *EventStore) RegisterEvent(eventType string, eventReflectType reflect.Type) {
	store.eventRegistry.Register(eventType, eventReflectType)
}
