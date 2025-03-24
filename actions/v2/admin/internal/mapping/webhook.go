package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/notifications"
	"github.com/samber/lo"
)

// MapToModelsWebhooks converts a slice of ModelWebhook to ModelsWebhooks
func MapToModelsWebhooks(webhooks []notifications.ModelWebhook) api.ModelsWebhooks {
	if webhooks == nil {
		return nil
	}

	return lo.Map(webhooks, MapToModelsWebhook)
}

// MapToModelsWebhook converts a single ModelWebhook to ModelsWebhook
func MapToModelsWebhook(w notifications.ModelWebhook, _ int) api.ModelsWebhook {
	if w == nil {
		return api.ModelsWebhook{}
	}

	return api.ModelsWebhook{
		Url:    w.GetURL(),
		Banned: w.Banned(),
	}
}
