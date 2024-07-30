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
	Model common.OldModel

	// ID is a draft transaction id.
	ID string `json:"id" example:"b356f7fa00cd3f20cce6c21d704cd13e871d28d714a5ebd0532f5a0e0cde63f7"`
	// Hex is a draft transaction hex.
	Hex string `json:"hex" example:"0100000002..."`
	// XpubID is a draft transaction's xpub used to sign transaction.
	XpubID string `json:"xpub_id" example:"bb8593f85ef8056a77026ad415f02128f3768906de53e9e8bf8749fe2d66cf50"`
	// ExpiresAt is a time when draft transaction expired.
	ExpiresAt time.Time `json:"expires_at" example:"2024-02-26T11:00:28.069911Z"`
	// Configuration contains draft transaction configuration.
	Configuration TransactionConfig `json:"configuration"`
	// Status is a draft transaction lastly monitored status.
	Status string `json:"status" example:"complete"`
	// FinalTxID is a final transaction id.
	FinalTxID string `json:"final_tx_id" example:"cfe30797f0b5fc098b32194e857569a7a1edd829fddf3df4567796b738de386d"`
}
