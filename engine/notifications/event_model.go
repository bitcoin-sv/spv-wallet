package notifications

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// RawEvent - event type
type RawEvent struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

type StringEvent struct {
	Value string
}

type NumericEvent struct {
	Numeric int
}

type Events interface {
	StringEvent | NumericEvent
}

func InstantinateEvent[EventType Events]() *EventType {
	base := *new(EventType)
	return &base
}

func GetEventNameByType[EventType Events]() string {
	content := InstantinateEvent[EventType]()
	return reflect.TypeOf(content).Elem().Name()
}

func GetEventName[EventType Events](instance *EventType) string {
	return reflect.TypeOf(instance).Elem().Name()
}

func GetEventContent[EventType Events](raw *RawEvent) (*EventType, error) {
	model := InstantinateEvent[EventType]()
	if raw.Type != GetEventName[EventType](model) {
		return nil, fmt.Errorf("Wrong type")
	}

	if err := json.Unmarshal(raw.Content, &model); err != nil {
		return nil, errors.Wrap(err, "Cannot unmarshall the content json")
	}
	return model, nil
}

func NewRawEvent[EventType Events](namedEvent *EventType) *RawEvent {
	asJson, _ := json.Marshal(namedEvent)
	return &RawEvent{
		Type:    GetEventName[EventType](namedEvent),
		Content: asJson,
	}
}
