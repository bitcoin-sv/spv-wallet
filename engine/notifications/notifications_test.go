package notifications

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockNotificationFrame struct {
	Value int `json:"value"`
}

func (mockNotificationFrame) GetType() string {
	return "mock-notification"
}

func newMockEvent(value int) *RawEvent {
	return NewRawEvent(mockNotificationFrame{
		Value: value,
	})
}

type mockNotifier struct {
	delay   *time.Duration
	channel chan *RawEvent
	output  []*RawEvent
}

func (m *mockNotifier) consumer(ctx context.Context) {
	for {
		select {
		case event := <-m.channel:
			m.output = append(m.output, event)
			if m.delay != nil {
				sleepWithContext(ctx, *m.delay)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *mockNotifier) assertOutput(t *testing.T, expected []int) {
	assert.Equal(t, len(expected), len(m.output))
	if len(expected) == len(m.output) {
		for i := 0; i < len(expected); i++ {
			actualEvent, err := GetEventContent[mockNotificationFrame](m.output[i])
			assert.NoError(t, err)
			assert.Equal(t, expected[i], actualEvent.Value)
		}
	}
}

func newMockNotifier(ctx context.Context, chanLength int) *mockNotifier {
	notifier := &mockNotifier{
		channel: make(chan *RawEvent, chanLength),
	}

	go notifier.consumer(ctx)
	return notifier
}

func sleepWithContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
	case <-timer.C:
	}
}

func TestNotifications(t *testing.T) {
	t.Run("one notifier", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx)
		notifier := newMockNotifier(ctx, 100)
		n.AddNotifier("test", notifier.channel)

		expected := []int{}
		for i := 0; i < 10; i++ {
			n.Notify(newMockEvent(i))
			expected = append(expected, i)
		}

		time.Sleep(100 * time.Millisecond)
		cancel()

		notifier.assertOutput(t, expected)
	})

	t.Run("two notifiers", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx)
		notifier1 := newMockNotifier(ctx, 100)
		notifier2 := newMockNotifier(ctx, 100)
		n.AddNotifier("notifier1", notifier1.channel)
		n.AddNotifier("notifier2", notifier2.channel)

		expected := []int{}
		for i := 0; i < 10; i++ {
			n.Notify(newMockEvent(i))
			expected = append(expected, i)
		}

		time.Sleep(100 * time.Millisecond)
		cancel()

		notifier1.assertOutput(t, expected)
		notifier2.assertOutput(t, expected)
	})

	t.Run("more notifications than output chan buffer", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx)
		outputChanLength := 10
		numberOfEvents := 50 // 50 > 10
		notifier := newMockNotifier(ctx, outputChanLength)
		n.AddNotifier("test", notifier.channel)

		expected := []int{}
		for i := 0; i < numberOfEvents; i++ {
			n.Notify(newMockEvent(i))
			// we have to delay of putting new events because the output chan buffer will not contain all of the events in its buffer
			// so, this way, we let consumer to pop events from the queue
			time.Sleep(1 * time.Millisecond)
			expected = append(expected, i)
		}

		time.Sleep(500 * time.Millisecond)
		cancel()

		notifier.assertOutput(t, expected)
	})

	t.Run("slow and fast consumers", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx)
		outputChanLength := 10
		numberOfEvents := 50 // 50 > 10

		notifier1 := newMockNotifier(ctx, outputChanLength)
		veryLongDelay := 1 * time.Hour // it means that notifier will pop only one event
		notifier1.delay = &veryLongDelay

		notifier2 := newMockNotifier(ctx, outputChanLength)
		n.AddNotifier("notifier1", notifier1.channel)
		n.AddNotifier("notifier2", notifier2.channel)

		expected := []int{}
		for i := 0; i < numberOfEvents; i++ {
			n.Notify(newMockEvent(i))
			time.Sleep(1 * time.Millisecond)
			expected = append(expected, i)
		}

		time.Sleep(500 * time.Millisecond)
		cancel()

		// even though notifier1 is slow, it should not block notifier2
		notifier1.assertOutput(t, []int{0})
		notifier2.assertOutput(t, expected)
	})

	t.Run("buffered and unbuffered channels", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx)
		outputChanLength := 10
		numberOfEvents := 50 // 50 > 10

		notifier1 := newMockNotifier(ctx, 1)
		notifier2 := newMockNotifier(ctx, outputChanLength)
		n.AddNotifier("notifier1", notifier1.channel)
		n.AddNotifier("notifier2", notifier2.channel)

		expected := []int{}
		for i := 0; i < numberOfEvents; i++ {
			n.Notify(newMockEvent(i))
			time.Sleep(1 * time.Millisecond)
			expected = append(expected, i)
		}

		time.Sleep(500 * time.Millisecond)
		cancel()

		notifier1.assertOutput(t, expected)
		notifier2.assertOutput(t, expected)
	})
}
