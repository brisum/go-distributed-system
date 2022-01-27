package event

import (
	"distributes_system/lib/event_sourcing"
)

type AccountCreatedEvent struct {
	firstName string
	lastName  string
}

func NewAccountCreatedEvent(firstName string, lastName string) *AccountCreatedEvent {
	return &AccountCreatedEvent{
		firstName: firstName,
		lastName:  lastName,
	}
}

func NewEmptyAccountCreatedEvent() *AccountCreatedEvent {
	return &AccountCreatedEvent{}
}

func (event *AccountCreatedEvent) GetFirstName() string {
	return event.firstName
}

func (event *AccountCreatedEvent) GetLastName() string {
	return event.lastName
}

func (event *AccountCreatedEvent) ToDataStorage() *eventsourcing.DataStorage {
	storage := eventsourcing.NewEmptyDataStorage()

	storage.Set("firstName", event.firstName)
	storage.Set("lastName", event.lastName)
	storage.Set("balance/cash", 0)
	storage.Set("balance/bonus", 0)

	return storage
}

func (event *AccountCreatedEvent) FromDataStorage(storage eventsourcing.DataStorage) {
	event.firstName = storage.Get("firstName").(string)
	event.lastName = storage.Get("lastName").(string)
}
