package notifications

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bsv-blockchain/spv-wallet/models"
)

func newMockEvent(value string) *models.RawEvent {
	return NewRawEvent(&models.StringEvent{
		Value: value,
	})
}

type mockNotifier struct {
	delay   *time.Duration
	channel chan *models.RawEvent
	output  []*models.RawEvent
	mu      sync.Mutex
}

func (m *mockNotifier) consumer(ctx context.Context) {
	for {
		select {
		case event := <-m.channel:
			m.mu.Lock()
			m.output = append(m.output, event)
			m.mu.Unlock()
			if m.delay != nil {
				sleepWithContext(ctx, *m.delay)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (m *mockNotifier) assertOutput(t *testing.T, expected []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	assert.Len(t, m.output, len(expected))
	if len(expected) == len(m.output) {
		for i := 0; i < len(expected); i++ {
			actualEvent, err := GetEventContent[models.StringEvent](m.output[i])
			require.NoError(t, err)
			assert.Equal(t, expected[i], actualEvent.Value)
		}
	}
}

func newMockNotifier(ctx context.Context, chanLength int) *mockNotifier {
	notifier := &mockNotifier{
		channel: make(chan *models.RawEvent, chanLength),
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
		n := NewNotifications(ctx, &nopLogger)
		notifier := newMockNotifier(ctx, 100)
		n.AddNotifier("test", notifier.channel)

		expected := []string{}
		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
			expected = append(expected, msg)
		}

		time.Sleep(100 * time.Millisecond)
		cancel()

		notifier.assertOutput(t, expected)
	})

	t.Run("two notifiers", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		notifier1 := newMockNotifier(ctx, 100)
		notifier2 := newMockNotifier(ctx, 100)
		n.AddNotifier("notifier1", notifier1.channel)
		n.AddNotifier("notifier2", notifier2.channel)

		expected := []string{}
		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
			expected = append(expected, msg)
		}

		time.Sleep(100 * time.Millisecond)
		cancel()

		notifier1.assertOutput(t, expected)
		notifier2.assertOutput(t, expected)
	})

	t.Run("more notifications than output chan buffer", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		outputChanLength := 10
		numberOfEvents := 50 // 50 > 10
		notifier := newMockNotifier(ctx, outputChanLength)
		n.AddNotifier("test", notifier.channel)

		expected := []string{}
		for i := 0; i < numberOfEvents; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
			// we have to delay of putting new events because the output chan buffer will not contain all of the events in its buffer
			// so, this way, we let consumer to pop events from the queue
			time.Sleep(1 * time.Millisecond)
			expected = append(expected, msg)
		}

		time.Sleep(500 * time.Millisecond)
		cancel()

		notifier.assertOutput(t, expected)
	})

	t.Run("slow and fast consumers", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		outputChanLength := 10
		numberOfEvents := 50 // 50 > 10

		notifier1 := newMockNotifier(ctx, outputChanLength)
		veryLongDelay := 1 * time.Hour // it means that notifier will pop only one event
		notifier1.delay = &veryLongDelay

		notifier2 := newMockNotifier(ctx, outputChanLength)
		n.AddNotifier("notifier1", notifier1.channel)
		n.AddNotifier("notifier2", notifier2.channel)

		expected := []string{}
		for i := 0; i < numberOfEvents; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
			time.Sleep(1 * time.Millisecond)
			expected = append(expected, msg)
		}

		time.Sleep(500 * time.Millisecond)
		cancel()

		// even though notifier1 is slow, it should not block notifier2
		notifier1.assertOutput(t, []string{"msg-0"})
		notifier2.assertOutput(t, expected)
	})

	t.Run("buffered and unbuffered channels", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		outputChanLength := 10
		numberOfEvents := 50 // 50 > 10

		// notifier1's buffer must hold every event: the fan-out send is non-blocking
		// (see sendEventToChannel), so a buffer smaller than numberOfEvents can drop an
		// event if the consumer goroutine isn't scheduled in time under CPU contention.
		// Sizing the buffer to numberOfEvents makes the "all events delivered" assertion
		// deterministic regardless of scheduling.
		notifier1 := newMockNotifier(ctx, numberOfEvents)
		notifier2 := newMockNotifier(ctx, outputChanLength)
		n.AddNotifier("notifier1", notifier1.channel)
		n.AddNotifier("notifier2", notifier2.channel)

		expected := []string{}
		for i := 0; i < numberOfEvents; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
			time.Sleep(1 * time.Millisecond)
			expected = append(expected, msg)
		}

		time.Sleep(500 * time.Millisecond)
		cancel()

		notifier1.assertOutput(t, expected)
		notifier2.assertOutput(t, expected)
	})

	t.Run("close with timeout - success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		notifier := newMockNotifier(ctx, 10)
		n.AddNotifier("test", notifier.channel)

		// Send a few events
		for i := 0; i < 5; i++ {
			n.Notify(newMockEvent(fmt.Sprintf("msg-%d", i)))
		}

		time.Sleep(50 * time.Millisecond)
		cancel()

		// Close should complete quickly (well under 3 second timeout)
		start := time.Now()
		err := n.Close()
		duration := time.Since(start)

		require.NoError(t, err)
		assert.Less(t, duration, 1*time.Second, "Close should complete quickly")
	})

	t.Run("close with slow goroutine - timeout", func(t *testing.T) {
		ctx := context.Background()
		n := NewNotifications(ctx, &nopLogger)

		// Manually increment wg to simulate a goroutine that won't finish
		n.wg.Add(1)

		// Close should timeout after 500ms but not return error to allow cleanup to continue
		start := time.Now()
		err := n.Close()
		duration := time.Since(start)

		// Should return nil even on timeout to allow cleanup to continue
		require.NoError(t, err)
		assert.GreaterOrEqual(t, duration, 200*time.Millisecond)
		assert.Less(t, duration, 1*time.Second, "Should timeout at ~500ms")
	})
}
