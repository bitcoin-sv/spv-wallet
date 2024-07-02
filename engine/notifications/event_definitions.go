package notifications

import "encoding/json"

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

	TransactionID string `json:"transactionId"`
	Status        string `json:"status"`
}

// Events - interface for all supported events
type Events interface {
	StringEvent | TransactionEvent
}
