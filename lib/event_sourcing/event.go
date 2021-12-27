package eventsourcing

import (
	"reflect"
	"regexp"
)

func GetEventType(event EventInterface) string {
	fullEventType := reflect.TypeOf(event).String()
	eventType := regexp.MustCompile(`^.*\.([a-zA-Z]+)Event$`).ReplaceAllString(fullEventType, "$1")

	return eventType
}
