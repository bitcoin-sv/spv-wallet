package beef

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
)

// TxQueryResult represents the result of a transaction query.
type TxQueryResult struct {
	SourceTXID string  // Source Transaction ID.
	RawHex     *string // Raw transaction in hexadecimal format.
	BeefHex    *string // BEEF-formatted transaction.
}

// IsBeef checks if the TxQueryResult contains a BEEF-formatted transaction.
func (tx TxQueryResult) IsBeef() bool { return tx.BeefHex != nil }

// IsRawTx checks if the TxQueryResult contains a raw transaction.
func (tx TxQueryResult) IsRawTx() bool { return tx.RawHex != nil }

type Repository interface {
	QueryTransactionInputParents(ctx context.Context, sourceTXIDs ...string) (TxQueryResultSlice, error)
}

type Service struct {
	repository Repository
}

func (s *Service) PrepareBEEF(tx *sdk.Transaction) (string, error) {
	return "", nil
}
