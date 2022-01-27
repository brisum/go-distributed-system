package domain

import (
	eventsourcing "distributes_system/lib/event_sourcing"
	accountDomainEvent "distributes_system/project/virtual_pay_network/domain/account/domain/event"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"reflect"
)

type balance struct {
	cash  int
	bonus int
}

type AccountAggregate struct {
	aggregateRoot eventsourcing.AggregateRoot
	firstName     string
	lastName      string
	balance       balance
}

func NewAccountAggregate(entityUuid uuid.UUID) *AccountAggregate {
	aggregate := AccountAggregate{
		aggregateRoot: *eventsourcing.NewAggregateRoot("Account", entityUuid, 2),
	}

	aggregate.firstName = ""
	aggregate.lastName = ""
	aggregate.balance = balance{cash: 0, bonus: 0}

	return &aggregate
}

func (aggregate *AccountAggregate) GetEntityType() string {
	return aggregate.aggregateRoot.GetEntityType()
}

func (aggregate *AccountAggregate) GetEntityUuid() uuid.UUID {
	return aggregate.aggregateRoot.GetEntityUuid()
}

func (aggregate *AccountAggregate) SetVersion(version int) {
	aggregate.aggregateRoot.SetVersion(version)
}

func (aggregate *AccountAggregate) GetVersion() int {
	return aggregate.aggregateRoot.GetVersion()
}

func (aggregate *AccountAggregate) GetSnapshotStrategy() int {
	return aggregate.aggregateRoot.GetSnapshotStrategy()
}

func (aggregate *AccountAggregate) GetEvents() *eventsourcing.EventStream {
	return aggregate.aggregateRoot.GetUncommittedEvents()
}

func (aggregate *AccountAggregate) ProcessEvent(event eventsourcing.EventInterface) {
	eventType := eventsourcing.GetEventType(event)
	methodName := "Apply" + eventType + "Event"
	methodArguments := []reflect.Value{reflect.ValueOf(event)}

	reflect.ValueOf(aggregate).MethodByName(methodName).Call(methodArguments)
	aggregate.aggregateRoot.AppendEvent(event)
}

func (aggregate *AccountAggregate) ApplyAccountCreatedEvent(event eventsourcing.EventInterface) {
	castedEvent := event.(*accountDomainEvent.AccountCreatedEvent)
	aggregate.firstName = castedEvent.GetFirstName()
	aggregate.lastName = castedEvent.GetLastName()
}

func (aggregate *AccountAggregate) ApplyBalanceIncreasedEvent(event eventsourcing.EventInterface) {
	castedEvent := event.(*accountDomainEvent.BalanceIncreasedEvent)
	aggregate.balance.cash = aggregate.balance.cash + castedEvent.GetCash()
	aggregate.balance.bonus = aggregate.balance.bonus + castedEvent.GetBonus()
}

func (aggregate *AccountAggregate) ToDataStorage() *eventsourcing.DataStorage {
	storage := eventsourcing.NewEmptyDataStorage()

	storage.Set("firstName", aggregate.firstName)
	storage.Set("lastName", aggregate.lastName)
	storage.Set("balance/cash", aggregate.balance.cash)
	storage.Set("balance/bonus", aggregate.balance.bonus)

	return storage
}

func (aggregate *AccountAggregate) FromDataStorage(storage eventsourcing.DataStorage) {
	aggregate.firstName = storage.Get("firstName").(string)
	aggregate.lastName = storage.Get("lastName").(string)
	aggregate.balance.cash = int(storage.Get("balance/cash").(float64))
	aggregate.balance.bonus = int(storage.Get("balance/bonus").(float64))
}

func (aggregate *AccountAggregate) CreateEventFromDataStorage(
	eventType string,
	storage eventsourcing.DataStorage,
) (eventsourcing.EventInterface, error) {
	switch eventType {
	case "AccountCreated":
		event := accountDomainEvent.NewAccountCreatedEvent("", "")
		event.FromDataStorage(storage)
		return event, nil

	case "BalanceIncreased":
		event := accountDomainEvent.NewBalanceIncreasedEvent(0, 0)
		event.FromDataStorage(storage)
		return event, nil
	}

	return nil, errors.New(fmt.Sprintf("Event \"%s\" not created.", eventType))
}
