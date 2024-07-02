package notifications

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventParsing(t *testing.T) {
	t.Run("parse the raw event to actual event type", func(t *testing.T) {
		source := NewRawEvent(&StringEvent{
			Value: "1",
		})
		asJSON, _ := json.Marshal(source)

		var target RawEvent
		_ = json.Unmarshal(asJSON, &target)
		assert.Equal(t, source.Type, target.Type)

		actualEvent, err := GetEventContent[StringEvent](&target)
		assert.NoError(t, err)
		assert.Equal(t, "1", actualEvent.Value)
	})

	t.Run("event name", func(t *testing.T) {
		assert.Equal(t, "StringEvent", GetEventNameByType[StringEvent]())
		var numericEventInstance *StringEvent
		assert.Equal(t, "StringEvent", GetEventName(numericEventInstance))
	})
}
