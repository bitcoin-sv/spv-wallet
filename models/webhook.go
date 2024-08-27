package models

// Webhook is a webhook model
// TokenHeader and TokenValue are not exposed because of security reasons
type Webhook struct {
	URL    string `json:"url"`
	Banned bool   `json:"banned"`
}
