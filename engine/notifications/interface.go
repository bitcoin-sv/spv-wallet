package notifications

import (
	"context"
	"time"
)

// WebhooksRepository is an interface for managing webhooks.
type WebhooksRepository interface {
	CreateWebhook(ctx context.Context, url, tokenHeader, tokenValue string) error
	RemoveWebhook(ctx context.Context, url string) error
	BanWebhook(ctx context.Context, url string, bannedTo time.Time) error
	GetWebhooks(ctx context.Context) ([]*WebhookModel, error)
}
