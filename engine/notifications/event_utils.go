package notifications

import (
	"encoding/json"
	"reflect"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// InstantinateEvent creates a new instance of the event type passed as a type parameter.
func InstantinateEvent[EventType models.Events]() *EventType {
	base := *new(EventType)
	return &base
}

// GetEventNameByType returns the name of the event type passed as a type parameter.
func GetEventNameByType[EventType models.Events]() string {
	content := InstantinateEvent[EventType]()
	return reflect.TypeOf(content).Elem().Name()
}

// GetEventName returns the name of the event type passed as a parameter.
func GetEventName[EventType models.Events](instance *EventType) string {
	return reflect.TypeOf(instance).Elem().Name()
}

// GetEventContent returns the content of the raw event passed as a parameter.
func GetEventContent[EventType models.Events](raw *models.RawEvent) (*EventType, error) {
	model := InstantinateEvent[EventType]()
	if raw.Type != GetEventName(model) {
		return nil, spverrors.Newf("wrong type")
	}

	if err := json.Unmarshal(raw.Content, &model); err != nil {
		return nil, spverrors.Wrapf(err, "cannot unmarshall the content json")
	}
	return model, nil
}

// NewRawEvent creates a new raw event from actual event object.
func NewRawEvent[EventType models.Events](namedEvent *EventType) *models.RawEvent {
	asJSON, _ := json.Marshal(namedEvent)
	return &models.RawEvent{
		Type:    GetEventName(namedEvent),
		Content: asJSON,
	}
}

// Notify is a utility generc function which allows to push a new event to the notification system.
func Notify[EventType models.Events](n *Notifications, event *EventType) {
	rawEvent := NewRawEvent(event)
	n.Notify(rawEvent)
}
