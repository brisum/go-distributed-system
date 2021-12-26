package eventsourcing

import "distributes_system/lib/datastorage"

type EventInterface interface {
	ToDataStorage() *datastorage.DataStorage
	FromDataStorage(dataStorage datastorage.DataStorage)
}
