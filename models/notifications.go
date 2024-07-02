package models

// SubscribeRequestBody represents the request body for the subscribe endpoint.
type SubscribeRequestBody struct {
	URL         string `json:"url"`
	TokenHeader string `json:"tokenHeader"`
	TokenValue  string `json:"tokenValue"`
}

// UnsubscribeRequestBody represents the request body for the unsubscribe endpoint.
type UnsubscribeRequestBody struct {
	URL string `json:"url"`
}
