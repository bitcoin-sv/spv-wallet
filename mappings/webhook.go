package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToWebhookContract will map the webhook model from spv-wallet engine to the spv-wallet-models contract
func MapToWebhookContract(w notifications.ModelWebhook) *models.Webhook {
	if w == nil {
		return nil
	}

	return &models.Webhook{
		URL:    w.GetURL(),
		Banned: w.Banned(),
	}
}
