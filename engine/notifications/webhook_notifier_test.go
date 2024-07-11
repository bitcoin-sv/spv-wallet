package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	interceptor     func(req *http.Request) (*http.Response, error)
	url             string
	receivedBatches [][]*models.RawEvent
}

var nopLogger = zerolog.Nop()

func newMockClient(url string) *mockClient {
	mc := &mockClient{
		receivedBatches: make([][]*models.RawEvent, 0),
		url:             url,
	}

	customResponder := func(req *http.Request) (*http.Response, error) {
		if mc.interceptor != nil {
			res, err := mc.interceptor(req)
			if res != nil {
				return res, err
			}
		}

		// Read the body from the incoming request
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), err
		}

		var events []*models.RawEvent
		err = json.Unmarshal(body, &events)
		if err != nil {
			return httpmock.NewStringResponse(500, ""), err
		}

		mc.receivedBatches = append(mc.receivedBatches, events)

		return httpmock.NewStringResponse(200, "OK"), nil
	}

	httpmock.RegisterResponder("POST", url, customResponder)

	return mc
}

func (mc *mockClient) assertEvents(t *testing.T, expected []string) {
	flatten := make([]*models.RawEvent, 0)
	for _, batch := range mc.receivedBatches {
		flatten = append(flatten, batch...)
	}
	assert.Equal(t, len(expected), len(flatten))
	if len(expected) == len(flatten) {
		for i := 0; i < len(expected); i++ {
			actualEvent, err := GetEventContent[models.StringEvent](flatten[i])
			assert.NoError(t, err)
			assert.Equal(t, expected[i], actualEvent.Value)
		}
	}
}

func (mc *mockClient) assertEventsWereSentInBatches(t *testing.T, expected bool) {
	result := false
	for _, batch := range mc.receivedBatches {
		if len(batch) > 1 {
			result = true
			break
		}
	}
	assert.Equal(t, expected, result)
}

type mockModelWebhook struct {
	BannedTo    *time.Time
	URL         string
	TokenHeader string
	TokenValue  string
	deleted     bool
}

func (m *mockModelWebhook) Banned() bool {
	return m.BannedTo != nil
}

func (m *mockModelWebhook) Deleted() bool {
	return m.deleted
}

func (m *mockModelWebhook) GetURL() string {
	return m.URL
}

func (m *mockModelWebhook) GetTokenHeader() string {
	return m.TokenHeader
}

func (m *mockModelWebhook) GetTokenValue() string {
	return m.TokenValue
}

func (m *mockModelWebhook) MarkUntil(bannedTo time.Time) {
	m.BannedTo = &bannedTo
}

func (m *mockModelWebhook) Refresh(tokenHeader, tokenValue string) {
	m.BannedTo = nil
	m.deleted = false
	m.TokenHeader = tokenHeader
	m.TokenValue = tokenValue
}

func newMockWebhookModel(url, tokenHeader, tokenValue string) *mockModelWebhook {
	return &mockModelWebhook{
		URL:         url,
		TokenHeader: tokenHeader,
		TokenValue:  tokenValue,
	}
}

func TestWebhookNotifier(t *testing.T) {
	t.Run("one webhook notifier", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")

		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		notifier := NewWebhookNotifier(ctx, &nopLogger, newMockWebhookModel(client.url, "", ""), make(chan string))
		n.AddNotifier(client.url, notifier.Channel)

		expected := []string{}
		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
			expected = append(expected, msg)
		}

		time.Sleep(100 * time.Millisecond)
		cancel()

		client.assertEvents(t, expected)
		client.assertEventsWereSentInBatches(t, true)
	})

	t.Run("two webhook notifiers", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client1 := newMockClient("http://localhost:8080")
		client2 := newMockClient("http://localhost:8081")

		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)

		notifier1 := NewWebhookNotifier(ctx, &nopLogger, newMockWebhookModel(client1.url, "", ""), make(chan string))
		n.AddNotifier(client1.url, notifier1.Channel)

		notifier2 := NewWebhookNotifier(ctx, &nopLogger, newMockWebhookModel(client2.url, "", ""), make(chan string))
		n.AddNotifier(client2.url, notifier2.Channel)

		expected := []string{}
		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
			expected = append(expected, msg)
		}

		time.Sleep(100 * time.Millisecond)
		cancel()

		client1.assertEvents(t, expected)
		client1.assertEventsWereSentInBatches(t, true)

		client2.assertEvents(t, expected)
		client2.assertEventsWereSentInBatches(t, true)
	})

	t.Run("no batches when notifications are put slowly", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")

		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		notifier := NewWebhookNotifier(ctx, &nopLogger, newMockWebhookModel(client.url, "", ""), make(chan string))
		n.AddNotifier(client.url, notifier.Channel)

		expected := []string{}
		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
			time.Sleep(100 * time.Microsecond)
			expected = append(expected, msg)
		}

		time.Sleep(100 * time.Millisecond)
		cancel()

		client.assertEvents(t, expected)
		client.assertEventsWereSentInBatches(t, false)
	})

	t.Run("with retries", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")
		k := 0
		client.interceptor = func(_ *http.Request) (*http.Response, error) {
			if k < 1 {
				k++
				return httpmock.NewStringResponse(408, ""), fmt.Errorf("Timeout")
			}
			return nil, nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		notifier := NewWebhookNotifier(ctx, &nopLogger, newMockWebhookModel(client.url, "", ""), make(chan string))
		n.AddNotifier(client.url, notifier.Channel)

		expected := []string{}
		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
			expected = append(expected, msg)
		}

		time.Sleep(1500 * time.Millisecond)
		cancel()

		client.assertEvents(t, expected)
	})

	t.Run("ban webhook", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")
		client.interceptor = func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(408, ""), fmt.Errorf("Timeout")
		}

		banMsg := make(chan string)
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		notifier := NewWebhookNotifier(ctx, &nopLogger, newMockWebhookModel(client.url, "", ""), banMsg)
		n.AddNotifier(client.url, notifier.Channel)

		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
		}

		banHasBeenTriggered := false
	loop:
		for {
			select {
			case <-time.After(3 * time.Second):
				break loop
			case url := <-banMsg:
				if url == client.url {
					banHasBeenTriggered = true
				}
			}
		}
		cancel()

		assert.Equal(t, true, banHasBeenTriggered)
	})

	t.Run("with token", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		tokenHeader := "test"
		tokenValue := "token"

		waitForCall := make(chan bool)
		client := newMockClient("http://localhost:8080")
		allGood := false
		client.interceptor = func(req *http.Request) (*http.Response, error) {
			defer func() {
				waitForCall <- true
			}()
			if tokenValue == req.Header.Get(tokenHeader) {
				allGood = true
			}
			return nil, nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx, &nopLogger)
		notifier := NewWebhookNotifier(ctx, &nopLogger, newMockWebhookModel(client.url, tokenHeader, tokenValue), make(chan string))
		n.AddNotifier(client.url, notifier.Channel)

		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("msg-%d", i)
			n.Notify(newMockEvent(msg))
		}

		<-waitForCall
		cancel()

		assert.Equal(t, true, allGood)
	})
}
