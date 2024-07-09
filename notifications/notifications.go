package notifications

import (
	"context"
	"sync"
)

// Event - event type
type Event any

const lengthOfInputChannel = 100

// Notifications - service for sending events to multiple notifiers
type Notifications struct {
	inputChannel   chan Event
	outputChannels *sync.Map //[string, chan Event]
}

// AddNotifier - add notifier by key
func (n *Notifications) AddNotifier(key string, ch chan Event) {
	n.outputChannels.Store(key, ch)
}

// RemoveNotifier - remove notifier by key
func (n *Notifications) RemoveNotifier(key string) {
	n.outputChannels.Delete(key)
}

// Notify - send event to all notifiers
func (n *Notifications) Notify(event Event) {
	n.inputChannel <- event
}

// exchange - exchange events between input and output channels, uses fan-out pattern
func (n *Notifications) exchange(ctx context.Context) {
	for {
		select {
		case event := <-n.inputChannel:
			n.outputChannels.Range(func(_, value any) bool {
				ch := value.(chan Event)
				n.sendEventToChannel(ch, event)
				return true
			})
		case <-ctx.Done():
			return
		}
	}
}

// sendEventToChannel - non blocking send event to channel
func (n *Notifications) sendEventToChannel(ch chan Event, event Event) {
	select {
	case ch <- event:
		// Successfully sent event
	default:
		// Channel is full, skip sending event
	}
}

// NewNotifications - creates a new instance of Notifications
func NewNotifications(ctx context.Context) *Notifications {
	n := &Notifications{
		inputChannel:   make(chan Event, lengthOfInputChannel),
		outputChannels: new(sync.Map),
	}

	go n.exchange(ctx)

	return n
}
