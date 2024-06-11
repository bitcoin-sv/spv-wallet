package chainstate

import (
	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/libsv/go-bc"
)

// TransactionInfo is the universal information about the transaction found from a chain provider
type TransactionInfo struct {
	BlockHash   string             `json:"block_hash,omitempty"` // mAPI
	BlockHeight int64              `json:"block_height"`         // mAPI
	ID          string             `json:"id"`                   // Transaction ID (Hex)
	Provider    string             `json:"provider,omitempty"`   // Provider is our internal source
	BUMP        *bc.BUMP           `json:"bump,omitempty"`       // Arc
	TxStatus    broadcast.TxStatus `json:"tx_status,omitempty"`  // Arc ONLY
}

// Valid validates TransactionInfo by checking if it contains
// BlockHash and MerkleProof (from mAPI) or BUMP (from Arc)
func (t *TransactionInfo) Valid() bool {
	arcInvalid := t.BUMP == nil
	invalid := t.BlockHash == "" || arcInvalid
	return !invalid
}
