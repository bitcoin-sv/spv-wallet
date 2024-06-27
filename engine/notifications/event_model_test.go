package notifications

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockEvent struct {
	Value int `json:"value"`
}

func (me MockEvent) GetType() string {
	return "mock-notification"
}

func TestEventParsing(t *testing.T) {
	t.Run("one notifier", func(t *testing.T) {
		source := NewRawEvent(&MockEvent{
			Value: 1,
		})
		asJSON, _ := json.Marshal(source)

		var target RawEvent
		_ = json.Unmarshal(asJSON, &target)
		assert.Equal(t, source.Type, target.Type)

		actualEvent, err := GetEventContent[MockEvent](&target)
		assert.NoError(t, err)
		assert.Equal(t, 1, actualEvent.Value)
	})
}
