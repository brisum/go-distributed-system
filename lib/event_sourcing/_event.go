package eventsourcing

//import (
//	"distributes_system/lib/datastorage"
//	uuid "github.com/satori/go.uuid"
//)
//
//type Event struct {
//	eventType      string
//	entityType     string
//	entityUuid     uuid.UUID
//	eventData      datastorage.DataStorage
//	promoter       string `default:""`
//	triggeredEvent string `default:""`
//}
//
//func NewEvent(
//	eventType string,
//	eventData datastorage.DataStorage,
//	promoter string,
//	triggeredEvent string) Event {
//
//	return Event{
//		eventType:      eventType,
//		eventData:      eventData,
//		promoter:       promoter,
//		triggeredEvent: triggeredEvent,
//	}
//}
//
//func (event *Event) GetEventType() string {
//	return event.eventType
//}
//
//func (event *Event) GetEntityType() string {
//	return event.entityType
//}
//
//func (event *Event) GetEventData() datastorage.DataStorage {
//	return event.eventData
//}
//
//func (event *Event) GetPromoter() string {
//	return event.promoter
//}
//
//func (event *Event) GetTriggeringEvent() string {
//	return event.triggeredEvent
//}
