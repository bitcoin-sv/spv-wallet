package notifications

import "time"

// WebhookModel - model for webhook stored in database
type WebhookModel struct {
	URL         string
	TokenHeader string
	Token       string
	CreatedAt   time.Time
}

func NewWebhookModel(url, tokenHeader, token string) *WebhookModel {
	return &WebhookModel{
		URL:         url,
		TokenHeader: tokenHeader,
		Token:       token,
		CreatedAt:   time.Now(),
	}
}
