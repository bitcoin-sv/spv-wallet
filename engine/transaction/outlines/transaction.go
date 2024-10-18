package outlines

import "github.com/bitcoin-sv/spv-wallet/engine/transaction"

// Transaction represents a transaction outline.
type Transaction struct {
	BEEF        string
	Annotations *transaction.Annotations
}
