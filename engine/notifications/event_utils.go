package notifications

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// InstantinateEvent creates a new instance of the event type passed as a type parameter.
func InstantinateEvent[EventType Events]() *EventType {
	base := *new(EventType)
	return &base
}

// GetEventNameByType returns the name of the event type passed as a type parameter.
func GetEventNameByType[EventType Events]() string {
	content := InstantinateEvent[EventType]()
	return reflect.TypeOf(content).Elem().Name()
}

// GetEventName returns the name of the event type passed as a parameter.
func GetEventName[EventType Events](instance *EventType) string {
	return reflect.TypeOf(instance).Elem().Name()
}

// GetEventContent returns the content of the raw event passed as a parameter.
func GetEventContent[EventType Events](raw *RawEvent) (*EventType, error) {
	model := InstantinateEvent[EventType]()
	if raw.Type != GetEventName(model) {
		return nil, fmt.Errorf("Wrong type")
	}

	if err := json.Unmarshal(raw.Content, &model); err != nil {
		return nil, errors.Wrap(err, "Cannot unmarshall the content json")
	}
	return model, nil
}

// NewRawEvent creates a new raw event from actual event object.
func NewRawEvent[EventType Events](namedEvent *EventType) *RawEvent {
	asJSON, _ := json.Marshal(namedEvent)
	return &RawEvent{
		Type:    GetEventName(namedEvent),
		Content: asJSON,
	}
}
