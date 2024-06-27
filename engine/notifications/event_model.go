package notifications

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// RawEvent - event type
type RawEvent struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

type EventContent interface {
	GetType() string
}

type GeneralPurposeEvent struct {
	Value string
}

func (GeneralPurposeEvent) GetType() string {
	return "general-purpose-event"
}

func GetEventContent[modelType EventContent](raw *RawEvent) (*modelType, error) {
	model := *new(modelType)
	if raw.Type != model.GetType() {
		return nil, fmt.Errorf("Wrong type")
	}

	if err := json.Unmarshal(raw.Content, &model); err != nil {
		return nil, errors.Wrap(err, "Cannot unmarshall the content json")
	}
	return &model, nil
}

func NewRawEvent(namedEvent EventContent) *RawEvent {
	asJson, _ := json.Marshal(namedEvent)
	return &RawEvent{
		Type:    namedEvent.GetType(),
		Content: asJson,
	}
}
