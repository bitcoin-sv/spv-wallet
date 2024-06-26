package notifications

import (
	"context"
	"time"
)

type WebhookModel struct {
	URL         string
	TokenHeader string
	TokenValue  string
	BannedTo    *time.Time
}

func (model *WebhookModel) Banned() bool {
	if model.BannedTo == nil {
		return false
	}
	ret := !time.Now().After(*model.BannedTo)
	return ret
}

type WebhooksRepository interface {
	CreateWebhook(ctx context.Context, url, tokenHeader, tokenValue string) error
	RemoveWebhook(ctx context.Context, url string) error
	BanWebhook(ctx context.Context, url string, bannedTo time.Time) error
	GetWebhooks(ctx context.Context) ([]*WebhookModel, error)
}
