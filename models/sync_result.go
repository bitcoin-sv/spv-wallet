package models

import "time"

// SyncResults is a model that represents a sync results.
type SyncResults struct {
	// LastMessage is a last message received during sync.
	LastMessage string `json:"last_message"`
	// Results is a slice of sync results.
	Results []*SyncResult `json:"results"`
}

// SyncResult is a model that represents a single sync result.
type SyncResult struct {
	// Action type broadcast, sync etc
	Action string `json:"action"`
	// ExecutedAt contains time when action was executed.
	ExecutedAt time.Time `json:"executed_at"`
	// Provider field is used for attempts(s).
	Provider string `json:"provider,omitempty"`
	// StatusMessage contains success or failure messages.
	StatusMessage string `json:"status_message"`
}
