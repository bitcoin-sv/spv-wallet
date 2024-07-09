package notifications

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

type mockRepository struct {
	webhooks []WebhookInterface
}

func (r *mockRepository) CreateWebhook(_ context.Context, url, tokenHeader, tokenValue string) error {
	r.webhooks = append(r.webhooks, newMockWebhookModel(url, tokenHeader, tokenValue))
	return nil
}

func (r *mockRepository) RemoveWebhook(_ context.Context, url string) error {
	for i, w := range r.webhooks {
		if w.GetURL() == url {
			r.webhooks = append(r.webhooks[:i], r.webhooks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("webhook not found")
}

func (r *mockRepository) GetWebhooks(_ context.Context) ([]WebhookInterface, error) {
	return r.webhooks, nil
}

func TestWebhookManager(t *testing.T) {
	t.Run("one webhook notifier previously subscribed", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")

		ctx, cancel := context.WithCancel(context.Background())

		n := NewNotifications(ctx)
		repo := &mockRepository{webhooks: []WebhookInterface{newMockWebhookModel(client.url, "", "")}}

		manager := NewWebhookManager(ctx, n, repo)
		time.Sleep(100 * time.Millisecond) // wait for manager to update notifiers
		defer manager.Stop()

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

	t.Run("one webhook notifier - subscribe", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")

		ctx, cancel := context.WithCancel(context.Background())

		n := NewNotifications(ctx)
		repo := &mockRepository{webhooks: []WebhookInterface{newMockWebhookModel(client.url, "", "")}}

		manager := NewWebhookManager(ctx, n, repo)
		time.Sleep(100 * time.Millisecond)
		defer manager.Stop()

		manager.Subscribe(ctx, client.url, "", "")
		time.Sleep(100 * time.Millisecond) // wait for manager to update notifiers

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
}
