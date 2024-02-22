package models

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/common"
)

// SyncTransaction is a model that represents a sync transaction specific fields.
type SyncTransaction struct {
	// Model is a common model that contains common fields for all models.
	common.Model

	// ID is a sync transaction id.
	ID string `json:"id"`
	// Configuration contains sync transaction configuration.
	Configuration SyncConfig `json:"configuration"`
	// LastAttempt contains last attempt time.
	LastAttempt time.Time `json:"last_attempt"`
	// Results contains sync transaction results.
	Results SyncResults `json:"results"`
	// BroadcastStatus contains broadcast status.
	BroadcastStatus string `json:"broadcast_status"`
	// P2PStatus contains p2p status.
	P2PStatus string `json:"p2p_status"`
	// SyncStatus contains sync status.
	SyncStatus string `json:"sync_status"`
}
