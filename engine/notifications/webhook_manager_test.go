package notifications

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

type mockRepository struct {
	webhooks []*WebhookModel
}

func (r *mockRepository) CreateWebhook(webhook *WebhookModel) error {
	r.webhooks = append(r.webhooks, webhook)
	return nil
}

func (r *mockRepository) RemoveWebhook(url string) error {
	for i, w := range r.webhooks {
		if w.URL == url {
			r.webhooks = append(r.webhooks[:i], r.webhooks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("webhook not found")
}

func (r *mockRepository) GetWebhooks() ([]*WebhookModel, error) {
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
		repo := &mockRepository{webhooks: []*WebhookModel{NewWebhookModel(client.url, "", "")}}

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
		repo := &mockRepository{webhooks: []*WebhookModel{NewWebhookModel(client.url, "", "")}}

		manager := NewWebhookManager(ctx, n, repo)
		time.Sleep(100 * time.Millisecond)
		defer manager.Stop()

		manager.Subscribe(&WebhookModel{URL: client.url})
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
