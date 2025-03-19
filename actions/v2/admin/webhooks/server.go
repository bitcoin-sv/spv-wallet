package webhooks

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	"github.com/rs/zerolog"
)

type webhooksService interface {
	SubscribeWebhook(ctx context.Context, url, tokenHeader, token string) error
	UnsubscribeWebhook(ctx context.Context, url string) error
	GetWebhooks(ctx context.Context) ([]notifications.ModelWebhook, error)
}

// APIAdminWebhooks represents server with admin API endpoints
type APIAdminWebhooks struct {
	webhooks webhooksService
	logger   *zerolog.Logger
}

// NewAPIAdminWebhooks creates a new APIAdminWebhooks
func NewAPIAdminWebhooks(webhooks webhooksService, log *zerolog.Logger) APIAdminWebhooks {
	logger := log.With().Str("api", "webhooks").Logger()

	return APIAdminWebhooks{
		webhooks: webhooks,
		logger:   &logger,
	}
}
