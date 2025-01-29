package outlines

import (
	"errors"
	"fmt"

	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
)

// TxQueryResult represents the result of a transaction query.
type TxQueryResult struct {
	SourceTXID string  // Source Transaction ID.
	RawHex     *string // Raw transaction in hexadecimal format.
	BeefHex    *string // BEEF-formatted transaction.
}

// IsBeef checks if the TxQueryResult contains a BEEF-formatted transaction.
func (tx TxQueryResult) IsBeef() bool {
	return tx.BeefHex != nil
}

// IsRawTx checks if the TxQueryResult contains a raw transaction.
func (tx TxQueryResult) IsRawTx() bool {
	return tx.RawHex != nil
}

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

// Value retrieves the SourceTx for a given ID or returns an empty SourceTx.
func (m SourceTxMap) Value(id string) SourceTx {
	return m[id]
}

// Add adds a TxQueryResult to the SourceTxMap after parsing it.
func (m SourceTxMap) Add(q *TxQueryResult) error {
	if q == nil {
		return nil
	}

	if q.IsBeef() && q.IsRawTx() {
		return ErrMutuallyExclusiveTxQueryResult
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

	return nil
}

// TxQueryResultSlice represents a slice of TxQueryResult.
type TxQueryResultSlice []*TxQueryResult

// SourceTxMap converts the slice to a SourceTxMap.
func (queryResults TxQueryResultSlice) SourceTxMap() (SourceTxMap, error) {
	sourceTXs := make(SourceTxMap)
	for _, q := range queryResults {
		if err := sourceTXs.Add(q); err != nil {
			return nil, fmt.Errorf("failed to add entry to source transaction map: %w", err)
		}
	}
	return sourceTXs, nil
}

// SourceTransactionBuilder builds transactions from TxQueryResults.
type SourceTransactionBuilder struct {
	Tx *sdk.Transaction // Root transaction.
}

// Build constructs the source transaction map and verifies SPV scripts.
func (b *SourceTransactionBuilder) Build(res TxQueryResultSlice) error {
	if b.Tx == nil {
		return ErrNilTxBuilder
	}

	sourceTxs, err := res.SourceTxMap()
	if err != nil {
		return fmt.Errorf("failed to convert query results to source transaction map: %w", err)
	}

	b.buildRecursive(b.Tx.Inputs, sourceTxs)

	for i, input := range b.Tx.Inputs {
		if input == nil || input.SourceTransaction == nil {
			return ErrInvalidTransactionInput
		}
		if _, err := spv.VerifyScripts(input.SourceTransaction); err != nil {
			return fmt.Errorf("SPV script verification failed for input %d: %w", i, err)
		}
	}
	return nil
}

// buildRecursive recursively attaches source transactions to inputs.
func (b *SourceTransactionBuilder) buildRecursive(inputs []*sdk.TransactionInput, sourceTxs SourceTxMap) {
	for _, input := range inputs {
		if input == nil {
			continue
		}

		sourceTxID := input.SourceTXID.String()
		val := sourceTxs.Value(sourceTxID)
		if val.IsZero() {
			continue
		}

		input.SourceTransaction = val.Tx
		if val.HadBeef {
			continue
		}

		b.buildRecursive(input.SourceTransaction.Inputs, sourceTxs)
	}
}

var (
	// ErrInvalidTransactionInput indicates that the SPV script verification failed
	// due to a nil transaction input or a missing source transaction.
	ErrInvalidTransactionInput = errors.New("SPV script verification failed: nil transaction input or missing source transaction")

	// ErrMutuallyExclusiveTxQueryResult indicates a conflict in the query result types.
	ErrMutuallyExclusiveTxQueryResult = errors.New("transaction query result must be either BEEF or RawTx type, not both")

	// ErrNilTxBuilder indicates a nil transaction in the builder.
	ErrNilTxBuilder = errors.New("transaction builder Tx field must be non-nil")
)
