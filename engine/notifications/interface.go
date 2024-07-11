package notifications

import (
	"context"
	"time"
)

// ModelWebhook is an interface for a webhook model.
type ModelWebhook interface {
	GetURL() string
	GetTokenHeader() string
	GetTokenValue() string
	MarkUntil(bannedTo time.Time)
	Refresh(tokenHeader, tokenValue string)
	Banned() bool
	Deleted() bool
}

// WebhooksRepository is an interface for managing webhooks.
type WebhooksRepository interface {
	Create(ctx context.Context, url, tokenHeader, tokenValue string) error
	Save(ctx context.Context, model ModelWebhook) error
	Delete(ctx context.Context, model ModelWebhook) error
	GetAll(ctx context.Context) ([]ModelWebhook, error)
	GetByURL(ctx context.Context, url string) (ModelWebhook, error)
}
