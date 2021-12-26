package domain

import (
	"distributes_system/lib/datastorage"
	eventsourcing "distributes_system/lib/event_sourcing"
	"distributes_system/project/virtual_pay_network/domain/account/domain/command"
	accountDomainEvent "distributes_system/project/virtual_pay_network/domain/account/domain/event"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"reflect"
)

type AccountAggregate struct {
	aggregateRoot eventsourcing.AggregateRoot
	firstName     string
	lastName      string
}

func NewAccountAggregate(entityUuid uuid.UUID) *AccountAggregate {
	aggregate := AccountAggregate{
		aggregateRoot: *eventsourcing.NewAggregateRoot("Account", entityUuid),
	}

	return &aggregate
}

func (aggregate *AccountAggregate) GetEntityType() string {
	return aggregate.aggregateRoot.GetEntityType()
}

func (aggregate *AccountAggregate) GetEntityUuid() uuid.UUID {
	return aggregate.aggregateRoot.GetEntityUuid()
}

func (aggregate *AccountAggregate) GetEvents() *eventsourcing.EventStream {
	return aggregate.aggregateRoot.GetUncommittedEvents()
}

func (aggregate *AccountAggregate) ProcessEvent(event eventsourcing.EventInterface) {

}

func (aggregate *AccountAggregate) CreateEventFromDataStorage(
	eventType string,
	storage datastorage.DataStorage,
) (eventsourcing.EventInterface, error) {
	switch eventType {
	case "AccountCreated":
		return accountDomainEvent.NewAccountCreatedEvent(
			storage.Get("firstName").(string),
			storage.Get("lastName").(string),
		), nil
	}

	return nil, errors.New(fmt.Sprintf("Event \"%s\" not created.", eventType))
}

func (aggregate *AccountAggregate) ProcessCreateAccountCommand(command command.CreateAccountCommand) {
	aggregate.aggregateRoot.AppendEvent(accountDomainEvent.NewAccountCreatedEvent(
		command.GetFirstName(),
		command.GetLastName(),
	))
}

func (aggregate *AccountAggregate) ApplyEvent(event eventsourcing.EventInterface) {
	eventType := eventsourcing.GetEventType(event)
	methodName := "Apply" + eventType + "Event"
	arguments := []reflect.Value{reflect.ValueOf(event)}

	reflect.ValueOf(aggregate).MethodByName(methodName).Call(arguments)
}

func (aggregate *AccountAggregate) ApplyAccountCreatedEvent(event eventsourcing.EventInterface) {
	castedEvent := event.(*accountDomainEvent.AccountCreatedEvent)
	aggregate.firstName = castedEvent.GetFirstName()
	aggregate.lastName = castedEvent.GetLastName()
}
