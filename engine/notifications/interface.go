package notifications

import (
	"context"
	"time"
)

type WebhooksRepository interface {
	CreateWebhook(ctx context.Context, url, tokenHeader, tokenValue string) error
	RemoveWebhook(ctx context.Context, url string) error
	BanWebhook(ctx context.Context, url string, bannedTo time.Time) error
	GetWebhooks(ctx context.Context) ([]*WebhookModel, error)
}
