package eventsourcing

type EventInterface interface {
	ToDataStorage() *DataStorage
	FromDataStorage(dataStorage DataStorage)
}
