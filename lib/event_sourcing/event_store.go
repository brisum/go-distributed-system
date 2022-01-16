package eventsourcing

import (
	"context"
	"distributes_system/lib/datastorage"
	"errors"
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
	tx, err := store.connection.Begin(*ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(*ctx)

	newVersion, err := store.saveEntity(ctx, tx, aggregate)
	if err != nil {
		return err
	}

	err = store.saveEvents(ctx, tx, aggregate)
	if err != nil {
		return err
	}

	err = store.saveSnapshot(ctx, tx, aggregate, newVersion)
	if err != nil {
		return err
	}

	err = tx.Commit(*ctx)
	if err != nil {
		return err
	}

	return nil
}

func (store *EventStore) saveEntity(ctx *context.Context, tx pgx.Tx, aggregate AggregateInterface) (int, error) {
	previousVersion := aggregate.GetVersion()
	newVersion := previousVersion + 1

	if 0 == previousVersion {
		result, err := tx.Exec(
			*ctx,
			"INSERT INTO entity (entity_type, entity_uuid, entity_version) "+
				"VALUES ($1, $2, $3)",
			aggregate.GetEntityType(),
			aggregate.GetEntityUuid().String(),
			newVersion,
		)
		if err != nil {
			return previousVersion, err
		}
		if 1 != result.RowsAffected() {
			return previousVersion, errors.New("Entity not inserted.")
		}

		return newVersion, nil
	}

	result, err := tx.Exec(
		*ctx,
		"UPDATE entity SET entity_version = $1 "+
			"WHERE entity_type = $2 AND entity_uuid = $3 AND entity_version = $4",
		newVersion,
		aggregate.GetEntityType(),
		aggregate.GetEntityUuid().String(),
		previousVersion,
	)
	if err != nil {
		return previousVersion, err
	}
	if 1 != result.RowsAffected() {
		return previousVersion, errors.New("Entity not updated.")
	}

	return newVersion, nil
}

func (store *EventStore) saveEvents(ctx *context.Context, tx pgx.Tx, aggregate AggregateInterface) error {
	stream := aggregate.GetEvents()

	reversedEvents := make([]EventInterface, 0)
	for _, event := range stream.GetEvents() {
		reversedEvents = append(reversedEvents, event)
	}

	for _, event := range reversedEvents {
		_, err := tx.Exec(
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

	return nil
}

func (store *EventStore) saveSnapshot(ctx *context.Context, tx pgx.Tx, aggregate AggregateInterface, newVersion int) error {
	return nil
}

//func (store *EventStore) RegisterEvent(eventType string, eventReflectType reflect.Type) {
//	store.eventRegistry.Register(eventType, eventReflectType)
//}

func (store *EventStore) Load(ctx *context.Context, aggregate AggregateInterface) error {
	err := store.loadSnapshot(ctx, aggregate)
	if err != nil {
		return err
	}

	err = store.loadEvents(ctx, aggregate)
	if err != nil {
		return err
	}

	err = store.loadEntity(ctx, aggregate)
	if err != nil {
		return err
	}

	return nil
}

func (store *EventStore) loadSnapshot(ctx *context.Context, aggregate AggregateInterface) error {
	return nil
}

func (store *EventStore) loadEvents(ctx *context.Context, aggregate AggregateInterface) error {
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

func (store *EventStore) loadEntity(ctx *context.Context, aggregate AggregateInterface) error {
	entityRow := EntityRow{}

	err := pgxscan.Get(
		*ctx,
		store.connection,
		&entityRow,
		`SELECT * FROM entity WHERE entity_type = $1 AND entity_uuid = $2 LIMIT 1`,
		aggregate.GetEntityType(),
		aggregate.GetEntityUuid().String(),
	)

	if err != nil {
		return err
	}

	aggregate.SetVersion(entityRow.EntityVersion)

	return nil
}
