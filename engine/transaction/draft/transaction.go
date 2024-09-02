package draft

import "github.com/bitcoin-sv/spv-wallet/engine/transaction"

// Transaction represents a transaction draft.
type Transaction struct {
	BEEF        string
	Annotations *transaction.Annotations
}
