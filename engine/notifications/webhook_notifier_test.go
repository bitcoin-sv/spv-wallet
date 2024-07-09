package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	interceptor     func(req *http.Request) (*http.Response, error)
	url             string
	receivedBatches [][]Event
}

func newMockClient(url string) *mockClient {
	mc := &mockClient{
		receivedBatches: make([][]Event, 0),
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

		var events []Event
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

func (mc *mockClient) assertEvents(t *testing.T, expected []Event) {
	flatten := make([]Event, 0)
	for _, batch := range mc.receivedBatches {
		flatten = append(flatten, batch...)
	}
	assert.Equal(t, len(expected), len(flatten))
	if len(expected) == len(flatten) {
		for i := 0; i < len(expected); i++ {
			assert.EqualValues(t, expected[i], flatten[i])
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

type mockWebhookModel struct {
	url         string
	tokenHeader string
	tokenValue  string
}

func (model *mockWebhookModel) GetURL() string {
	return model.url
}

func (model *mockWebhookModel) GetToken() (string, string) {
	return model.tokenHeader, model.tokenValue
}

func newMockWebhookModel(url, tokenHeader, tokenValue string) *mockWebhookModel {
	return &mockWebhookModel{
		url:         url,
		tokenHeader: tokenHeader,
		tokenValue:  tokenValue,
	}
}

func TestWebhookNotifier(t *testing.T) {
	t.Run("one webhook notifier", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")

		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx)
		notifier := NewWebhookNotifier(ctx, newMockWebhookModel(client.url, "", ""))
		n.AddNotifier(client.url, notifier.Channel)

		expected := []Event{}
		for i := 0; i < 10; i++ {
			n.Notify(i)
			expected = append(expected, i)
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
		n := NewNotifications(ctx)

		notifier1 := NewWebhookNotifier(ctx, newMockWebhookModel(client1.url, "", ""))
		n.AddNotifier(client1.url, notifier1.Channel)

		notifier2 := NewWebhookNotifier(ctx, newMockWebhookModel(client2.url, "", ""))
		n.AddNotifier(client2.url, notifier2.Channel)

		expected := []Event{}
		for i := 0; i < 10; i++ {
			n.Notify(i)
			expected = append(expected, i)
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
		n := NewNotifications(ctx)
		notifier := NewWebhookNotifier(ctx, newMockWebhookModel(client.url, "", ""))
		n.AddNotifier(client.url, notifier.Channel)

		expected := []Event{}
		for i := 0; i < 10; i++ {
			n.Notify(i)
			time.Sleep(100 * time.Microsecond)
			expected = append(expected, i)
		}

		time.Sleep(100 * time.Millisecond)
		cancel()

		client.assertEvents(t, expected)
		client.assertEventsWereSentInBatches(t, false)
	})

	t.Run("no batches when notifications are put slowly", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")

		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx)
		notifier := NewWebhookNotifier(ctx, newMockWebhookModel(client.url, "", ""))
		n.AddNotifier(client.url, notifier.Channel)

		expected := []Event{}
		for i := 0; i < 10; i++ {
			n.Notify(i)
			time.Sleep(100 * time.Microsecond)
			expected = append(expected, i)
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
		n := NewNotifications(ctx)
		notifier := NewWebhookNotifier(ctx, newMockWebhookModel(client.url, "", ""))
		n.AddNotifier(client.url, notifier.Channel)

		expected := []Event{}
		for i := 0; i < 10; i++ {
			n.Notify(i)
			expected = append(expected, i)
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

		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx)
		notifier := NewWebhookNotifier(ctx, newMockWebhookModel(client.url, "", ""))
		n.AddNotifier(client.url, notifier.Channel)

		for i := 0; i < 10; i++ {
			n.Notify(i)
		}

		time.Sleep(2500 * time.Millisecond)
		cancel()

		assert.Equal(t, true, notifier.banned())
	})

	t.Run("ban webhook from outside", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")
		ctx, cancel := context.WithCancel(context.Background())
		n := NewNotifications(ctx)
		notifier := NewWebhookNotifier(ctx, newMockWebhookModel(client.url, "", ""))
		n.AddNotifier(client.url, notifier.Channel)

		notifier.Ban()

		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, true, notifier.banned())

		for i := 0; i < 10; i++ {
			n.Notify(i)
		}

		time.Sleep(100 * time.Millisecond)
		cancel()

		client.assertEvents(t, []Event{})
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
		n := NewNotifications(ctx)
		notifier := NewWebhookNotifier(ctx, newMockWebhookModel(client.url, tokenHeader, tokenValue))
		n.AddNotifier(client.url, notifier.Channel)

		for i := 0; i < 10; i++ {
			n.Notify(i)
		}

		<-waitForCall
		cancel()

		assert.Equal(t, true, allGood)
	})
}
