package notifications

import "context"

type WebhookInterface interface {
	GetURL() string
	GetToken() (string, string) // key, value
}

type WebhooksRepository interface {
	CreateWebhook(ctx context.Context, url, tokenHeader, tokenValue string) error
	RemoveWebhook(ctx context.Context, url string) error
	GetWebhooks(ctx context.Context) ([]WebhookInterface, error)
}
