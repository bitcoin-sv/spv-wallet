package beef

import (
	"fmt"

	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
)

// SourceTx represents a source transaction.
type SourceTx struct {
	Tx      *sdk.Transaction // Parsed transaction.
	HadBeef bool             // Indicates if the transaction originated from a BEEF format.
}

// IsZero checks if the SourceTx is uninitialized.
func (s SourceTx) IsZero() bool {
	return s == SourceTx{}
}

// IsBeef checks if the SourceTx originated from a BEEF format.
func (s SourceTx) IsBeef() bool {
	return s.HadBeef
}

// SourceTxMap maps transaction IDs to SourceTx objects.
type SourceTxMap map[string]SourceTx

// Has checks if a transaction ID exists in the map.
func (m SourceTxMap) Has(id string) bool {
	_, ok := m[id]
	return ok
}

// Add adds a TxQueryResult to the SourceTxMap after parsing it.
func (m SourceTxMap) Add(q *TxQueryResult) error {
	if q == nil {
		return txerrors.ErrNilTxQueryResult
	}

	if q.IsBeef() {
		tx, err := sdk.NewTransactionFromBEEFHex(*q.BeefHex)
		if err != nil {
			return fmt.Errorf("failed to parse BEEF transaction: %w", err)
		}
		m[q.SourceTXID] = SourceTx{Tx: tx, HadBeef: true}
		return nil
	}

	if q.IsRawTx() {
		tx, err := sdk.NewTransactionFromHex(*q.RawHex)
		if err != nil {
			return fmt.Errorf("failed to parse raw transaction: %w", err)
		}
		m[q.SourceTXID] = SourceTx{Tx: tx, HadBeef: false}
		return nil
	}

	return txerrors.ErrTxQueryResultType
}

// Value retrieves the SourceTx for a given ID or returns an empty SourceTx.
func (m SourceTxMap) Value(id string) SourceTx {
	return m[id]
}

// TxQueryResultSlice represents a slice of TxQueryResult.
type TxQueryResultSlice []*TxQueryResult

// SourceTxMap converts the slice to a SourceTxMap.
func (qq TxQueryResultSlice) SourceTxMap() (SourceTxMap, error) {
	sourceTXs := make(SourceTxMap)
	for _, q := range qq {
		if err := sourceTXs.Add(q); err != nil {
			return nil, fmt.Errorf("failed to add entry to source transaction map: %w", err)
		}
	}
	return sourceTXs, nil
}

// SourceTransactionResolver set source transactions for each input given
// in the subject transaction based on the TxQueryResults.
type SourceTransactionResolver struct {
	subjectTx *sdk.Transaction
	sourceTxs SourceTxMap
}

// Resolve sets the source transaction per input of subject transaction and verifies SPV scripts.
func (s *SourceTransactionResolver) Resolve() error {
	if err := s.resolveRecursive(s.subjectTx.Inputs); err != nil {
		return fmt.Errorf("failed to resolve source transactions: %w", err)
	}

	for i, input := range s.subjectTx.Inputs {
		if input == nil || input.SourceTransaction == nil {
			return txerrors.ErrInvalidTransactionInput
		}
		if _, err := spv.VerifyScripts(input.SourceTransaction); err != nil {
			return fmt.Errorf("SPV script verification failed for input %d: %w", i, err)
		}
	}

	return nil
}

// resolveRecursive recursively attaches source transactions to inputs.
func (s *SourceTransactionResolver) resolveRecursive(inputs []*sdk.TransactionInput) error {
	for idx, input := range inputs {
		if input == nil {
			return txerrors.ErrNilTransactionInput
		}

		sourceTxID := input.SourceTXID.String()
		val := s.sourceTxs.Value(sourceTxID)
		if val.IsZero() {
			continue
		}

		input.SourceTransaction = val.Tx
		if val.HadBeef {
			continue
		}

		if err := s.resolveRecursive(input.SourceTransaction.Inputs); err != nil {
			return fmt.Errorf("Transaction %s failed to resolve source transaction for input %d: %w", sourceTxID, idx, err)
		}
	}

	return nil
}

// NewSourceTransactionResolver returns initialized source transaction resolver for given subject transaction
// and TxQueryResults.
func NewSourceTransactionResolver(tx *sdk.Transaction, slice TxQueryResultSlice) (*SourceTransactionResolver, error) {
	if tx == nil {
		return nil, txerrors.ErrNilSubjectTx
	}
	txs, err := slice.SourceTxMap()
	if err != nil {
		return nil, fmt.Errorf("failed to convert tx query result slice into map: %w", err)
	}

	return &SourceTransactionResolver{subjectTx: tx, sourceTxs: txs}, nil
}
