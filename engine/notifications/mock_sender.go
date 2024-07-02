package notifications

import (
	"context"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"
)

// StartSendingMockEvents - utility function to start sending some events in a predefined interval. It's useful for testing
func StartSendingMockEvents[EventType models.Events](ctx context.Context, notificationService *Notifications, duration time.Duration, prepare func(i int) *EventType) {
	i := 0
	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				notificationService.Notify(NewRawEvent(prepare(i)))
				i++
			}
		}
	}()
}
