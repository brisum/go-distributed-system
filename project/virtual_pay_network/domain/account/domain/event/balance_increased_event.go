package event

import (
	"distributes_system/lib/event_sourcing"
)

type BalanceIncreasedEvent struct {
	cash  int
	bonus int
}

func NewBalanceIncreasedEvent(cash int, bonus int) *BalanceIncreasedEvent {
	return &BalanceIncreasedEvent{
		cash:  cash,
		bonus: bonus,
	}
}

func NewEmptyBalanceIncreasedEvent() *BalanceIncreasedEvent {
	return &BalanceIncreasedEvent{}
}

func (event *BalanceIncreasedEvent) GetCash() int {
	return event.cash
}

func (event *BalanceIncreasedEvent) GetBonus() int {
	return event.bonus
}

func (event *BalanceIncreasedEvent) ToDataStorage() *eventsourcing.DataStorage {
	storage := eventsourcing.NewEmptyDataStorage()

	storage.Set("balance/cash", event.cash)
	storage.Set("balance/bonus", event.bonus)

	return storage
}

func (event *BalanceIncreasedEvent) FromDataStorage(storage eventsourcing.DataStorage) {
	event.cash = int(storage.Get("balance/cash").(float64))
	event.bonus = int(storage.Get("balance/bonus").(float64))
}
