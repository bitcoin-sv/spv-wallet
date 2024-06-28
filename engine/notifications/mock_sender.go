package notifications

import (
	"context"
	"time"
)

func StartSendingMockEvents[EventType Events](ctx context.Context, notificationService *Notifications, duration time.Duration, prepare func(i int) *EventType) {
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
