package notifications

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

type mockRepository struct {
	webhooks []ModelWebhook
}

func (r *mockRepository) Create(_ context.Context, url, tokenHeader, tokenValue string) error {
	model := newMockWebhookModel(url, tokenHeader, tokenValue)
	r.webhooks = append(r.webhooks, model)
	return nil
}

func (r *mockRepository) Save(_ context.Context, model ModelWebhook) error {
	for i, w := range r.webhooks {
		if w.GetURL() == model.GetURL() {
			r.webhooks[i] = model
			return nil
		}
	}
	r.webhooks = append(r.webhooks, model)
	return nil
}

func (r *mockRepository) Delete(_ context.Context, model ModelWebhook) error {
	for i, w := range r.webhooks {
		if w.GetURL() == model.GetURL() {
			webhook := r.webhooks[i].(*mockModelWebhook)
			webhook.deleted = true
			r.webhooks[i] = webhook
			return nil
		}
	}
	return nil
}

func (r *mockRepository) GetAll(_ context.Context) ([]ModelWebhook, error) {
	return r.webhooks, nil
}

func (r *mockRepository) GetByURL(_ context.Context, url string) (ModelWebhook, error) {
	for _, w := range r.webhooks {
		if w.GetURL() == url {
			return w, nil
		}
	}
	return nil, nil
}

func TestWebhookManager(t *testing.T) {
	t.Run("one webhook notifier previously subscribed", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")

		ctx, cancel := context.WithCancel(context.Background())

		n := NewNotifications(ctx, &nopLogger)
		repo := &mockRepository{webhooks: []ModelWebhook{newMockWebhookModel(client.url, "", "")}}

		manager := NewWebhookManager(ctx, &nopLogger, n, repo)
		time.Sleep(100 * time.Millisecond) // wait for manager to update notifiers
		defer manager.Stop()

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

	t.Run("one webhook notifier - subscribe", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")

		ctx, cancel := context.WithCancel(context.Background())

		n := NewNotifications(ctx, &nopLogger)
		repo := &mockRepository{webhooks: []ModelWebhook{newMockWebhookModel(client.url, "", "")}}

		manager := NewWebhookManager(ctx, &nopLogger, n, repo)
		time.Sleep(100 * time.Millisecond)
		defer manager.Stop()

		manager.Subscribe(ctx, client.url, "", "")
		time.Sleep(100 * time.Millisecond) // wait for manager to update notifiers

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
}
