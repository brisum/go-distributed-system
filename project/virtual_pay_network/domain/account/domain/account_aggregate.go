package domain

import (
	eventsourcing "distributes_system/lib/event_sourcing"
	"distributes_system/project/virtual_pay_network/domain/account/domain/command"
	"distributes_system/project/virtual_pay_network/domain/account/domain/event"
	uuid "github.com/satori/go.uuid"
)

type AccountAggregate struct {
	aggregateRoot eventsourcing.AggregateRoot
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

func (aggregate *AccountAggregate) ProcessCreateAccountCommand(command command.CreateAccountCommand) {
	aggregate.aggregateRoot.AppendEvent(event.NewAccountCreatedEvent(command.GetFirstName(), command.GetLastName()))
}

func (aggregate *AccountAggregate) ApplyAccountCreatedEvent() {

}
