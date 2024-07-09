package notifications

import (
	"encoding/json"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/assert"
)

func TestEventParsing(t *testing.T) {
	t.Run("parse the raw event to actual event type", func(t *testing.T) {
		source := NewRawEvent(&models.StringEvent{
			Value: "1",
		})
		asJSON, _ := json.Marshal(source)

		var target models.RawEvent
		_ = json.Unmarshal(asJSON, &target)
		assert.Equal(t, source.Type, target.Type)

		actualEvent, err := GetEventContent[models.StringEvent](&target)
		assert.NoError(t, err)
		assert.Equal(t, "1", actualEvent.Value)
	})

	t.Run("event name", func(t *testing.T) {
		assert.Equal(t, "StringEvent", GetEventNameByType[models.StringEvent]())
		var numericEventInstance *models.StringEvent
		assert.Equal(t, "StringEvent", GetEventName(numericEventInstance))
	})
}
