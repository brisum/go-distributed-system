package eventsourcing

import (
	"context"
	pgx "github.com/jackc/pgx/v4"
	"reflect"
)

type EventStore struct {
	connection *pgx.Conn
}

func NewEventStore(connection *pgx.Conn) *EventStore {
	store := EventStore{
		connection: connection,
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
			getEventType(event),
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

func getEventType(event EventInterface) string {
	return reflect.TypeOf(event).String()
}

//func (store *EventStore) Load(ctx *context.Context, aggregate *AggregateRootInterface) {
//
//}

//func (store *EventStore) _Load(ctx *context.Context, stream *EventStream) (*EventStream, error) {
//	eventRows := make([]EventRow, 0)
//
//	err := pgxscan.Select(
//		*ctx,
//		store.connection,
//		&eventRows,
//		`SELECT event_type, event_data FROM event WHERE entity_type = $1 and entity_uuid = $2`,
//		stream.GetEntityType(),
//		stream.GetEntityUuid().String(),
//	)
//
//	if err != nil {
//		return nil, err
//	}
//
//	for _, eventRow := range eventRows {
//		stream.AppendEvent(NewEvent(eventRow.EventType, eventRow.EventData, "", ""))
//	}
//
//	return stream, nil
//}
//
//func (store *EventStore) Append(ctx *context.Context, stream *EventStream) error {
//	reversedEvents := make([]Event, 0)
//	for _, event := range stream.GetEvents() {
//		reversedEvents = append(reversedEvents, event)
//	}
//
//	tx, err := store.connection.Begin(*ctx)
//	if err != nil {
//		return err
//	}
//	defer tx.Rollback(*ctx)
//
//	for _, event := range reversedEvents {
//		_, err = tx.Exec(
//			*ctx,
//			"INSERT INTO event (event_type, entity_type, entity_uuid, event_data, promoter, triggering_event) "+
//				"VALUES ($1, $2, $3, $4, $5, $6)",
//			event.GetEventType(),
//			stream.GetEntityType(),
//			stream.GetEntityUuid().String(),
//			event.GetEventData(),
//			event.GetPromoter(),
//			event.GetTriggeringEvent(),
//		)
//
//		if err != nil {
//			return err
//		}
//	}
//
//	err = tx.Commit(*ctx)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
