package models

import "encoding/json"

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

// RawEvent - the base event type
type RawEvent struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

// StringEvent - event with string value; can be used for generic messages and it's used for testing
type StringEvent struct {
	Value string `json:"value"`
}

// UserEvent - event with user identifier
type UserEvent struct {
	XPubID string `json:"xpubId"`
}

// TransactionEvent - event for transaction changes
type TransactionEvent struct {
	UserEvent `json:",inline"`

	TransactionID   string           `json:"transactionId"`
	Status          string           `json:"status"`
	XpubOutputValue map[string]int64 `json:"xpubOutputValue"`
	XpubInIDs	   []string          `json:"xpubInIds"`
}

// NOTICE: If you add a new event type, you must also update the Events interface

// Events - interface for all supported events
type Events interface {
	StringEvent | TransactionEvent
}
