package chainstate

import (
	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
)

// TransactionInfo is the universal information about the transaction found from a chain provider
type TransactionInfo struct {
	BlockHash   string             `json:"block_hash,omitempty"` // Block hash of the transaction
	BlockHeight int64              `json:"block_height"`         // Block height of the transaction
	ID          string             `json:"id"`                   // Transaction ID (Hex)
	Provider    string             `json:"provider,omitempty"`   // Provider is our internal source
	BUMP        *trx.MerklePath    `json:"bump,omitempty"`       // Merkle proof in BUMP format
	TxStatus    broadcast.TxStatus `json:"tx_status,omitempty"`  // Status of the transaction
}
