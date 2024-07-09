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

func (r *mockRepository) CreateWebhook(_ context.Context, url, tokenHeader, tokenValue string) error {
	model := newMockWebhookModel(url, tokenHeader, tokenValue)
	r.webhooks = append(r.webhooks, model)
	return nil
}

func (r *mockRepository) RemoveWebhook(_ context.Context, url string) error {
	for i, w := range r.webhooks {
		if w.URL == url {
			r.webhooks = append(r.webhooks[:i], r.webhooks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("webhook not found")
}

func (r *mockRepository) GetWebhooks(_ context.Context) ([]*WebhookModel, error) {
	return r.webhooks, nil
}

func (r *mockRepository) BanWebhook(_ context.Context, url string, banTime time.Time) error {
	for _, item := range r.webhooks {
		if item.URL == url {
			item.BannedTo = &banTime
		}
	}
	return nil
}

func TestWebhookManager(t *testing.T) {
	t.Run("one webhook notifier previously subscribed", func(t *testing.T) {
		httpmock.Reset()
		httpmock.Activate()
		defer httpmock.Deactivate()

		client := newMockClient("http://localhost:8080")

		ctx, cancel := context.WithCancel(context.Background())

		n := NewNotifications(ctx)
		repo := &mockRepository{webhooks: []*WebhookModel{newMockWebhookModel(client.url, "", "")}}

		manager := NewWebhookManager(ctx, n, repo)
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

		n := NewNotifications(ctx)
		repo := &mockRepository{webhooks: []*WebhookModel{newMockWebhookModel(client.url, "", "")}}

		manager := NewWebhookManager(ctx, n, repo)
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
