package notifications

import (
	"context"
	"time"
)

type ModelWebhook interface {
	GetURL() string
	GetTokenHeader() string
	GetTokenValue() string
	MarkDeleted()
	MarkBanned(bannedTo time.Time)
	Refresh(tokenHeader, tokenValue string)
	Banned() bool
}

// WebhooksRepository is an interface for managing webhooks.
type WebhooksRepository interface {
	Create(ctx context.Context, url, tokenHeader, tokenValue string) error
	Save(ctx context.Context, model ModelWebhook) error
	GetAll(ctx context.Context) ([]ModelWebhook, error)
	GetByURL(ctx context.Context, url string) (ModelWebhook, error)

	// CreateWebhook(ctx context.Context, url, tokenHeader, tokenValue string) error
	// RemoveWebhook(ctx context.Context, url string) error
	// BanWebhook(ctx context.Context, url string, bannedTo time.Time) error
	// GetWebhooks(ctx context.Context) ([]*WebhookModel, error)
}
