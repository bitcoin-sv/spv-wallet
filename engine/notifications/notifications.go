package notifications

import (
	"context"
	"sync"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/rs/zerolog"
)

const lengthOfInputChannel = 100

// Notifications - service for sending events to multiple notifiers
type Notifications struct {
	inputChannel   chan *models.RawEvent
	outputChannels *sync.Map //[string, chan *Event]
	burstLogger    *zerolog.Logger
}

// AddNotifier - add notifier by key
func (n *Notifications) AddNotifier(key string, ch chan *models.RawEvent) {
	n.outputChannels.Store(key, ch)
}

// RemoveNotifier - remove notifier by key
func (n *Notifications) RemoveNotifier(key string) {
	n.outputChannels.Delete(key)
}

// Notify - send event to all notifiers
func (n *Notifications) Notify(event *models.RawEvent) {
	n.inputChannel <- event
}

// exchange - exchange events between input and output channels, uses fan-out pattern
func (n *Notifications) exchange(ctx context.Context) {
	for {
		select {
		case event := <-n.inputChannel:
			n.outputChannels.Range(func(_, value any) bool {
				ch := value.(chan *models.RawEvent)
				n.sendEventToChannel(ch, event)
				return true
			})
		case <-ctx.Done():
			return
		}
	}
}

// sendEventToChannel - non blocking send event to channel
func (n *Notifications) sendEventToChannel(ch chan *models.RawEvent, event *models.RawEvent) {
	select {
	case ch <- event:
		// Successfully sent event
	default:
		n.burstLogger.Warn().Msg("Failed to send event to channel")
	}
}

// NewNotifications - creates a new instance of Notifications
func NewNotifications(ctx context.Context, parentLogger *zerolog.Logger) *Notifications {
	burstLogger := parentLogger.With().Logger().Sample(&zerolog.BurstSampler{
		Burst:  3,
		Period: 30 * time.Second,
	})
	n := &Notifications{
		inputChannel:   make(chan *models.RawEvent, lengthOfInputChannel),
		outputChannels: new(sync.Map),
		burstLogger:    &burstLogger,
	}

	go n.exchange(ctx)

	return n
}
