package models

import (
	"time"

	"github.com/bitcoin-sv/spv-wallet/models/common"
)

const (
	// DraftStatusDraft is when the transaction is a draft
	DraftStatusDraft string = "draft"

	// DraftStatusCanceled is when the draft is canceled
	DraftStatusCanceled string = "canceled"

	// DraftStatusExpired is when the draft has expired
	DraftStatusExpired string = "expired"

	// DraftStatusComplete is when the draft transaction is complete
	DraftStatusComplete string = "complete"
)

// DraftTransaction is a model that represents a draft transaction.
type DraftTransaction struct {
	// Model is a common model that contains common fields for all models.
	common.Model

	// ID is a draft transaction id.
	ID string `json:"id"`
	// Hex is a draft transaction hex.
	Hex string `json:"hex"`
	// XpubID is a draft transaction's xpub used to sign transaction.
	XpubID string `json:"xpub_id"`
	// ExpiresAt is a time when draft transaction expired.
	ExpiresAt time.Time `json:"expires_at"`
	// Configuration contains draft transaction configuration.
	Configuration TransactionConfig `json:"configuration"`
	// Status is a draft transaction lastly monitored status.
	Status string `json:"status"`
	// FinalTxID is a final transaction id.
	FinalTxID string `json:"final_tx_id"`
}
