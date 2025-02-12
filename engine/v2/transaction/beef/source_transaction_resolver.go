package beef

import (
	"errors"

	"github.com/bitcoin-sv/go-sdk/spv"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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
			return spverrors.Wrapf(errors.Join(err, txerrors.ErrInvalidBEEFHexInQueryResult), "failed to parse BEEF transaction")
		}
		m[q.SourceTXID] = SourceTx{Tx: tx, HadBeef: true}
		return nil
	}

	if q.IsRawTx() {
		tx, err := sdk.NewTransactionFromHex(*q.RawHex)
		if err != nil {
			return spverrors.Wrapf(errors.Join(err, txerrors.ErrInvalidRawHexInQueryResult), "failed to parse raw transaction")
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
			return nil, spverrors.Wrapf(err, "failed to add entry to source transaction map")
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
		return spverrors.Wrapf(err, "failed to resolve source transactions")
	}

	for i, input := range s.subjectTx.Inputs {
		if input == nil || input.SourceTransaction == nil {
			return txerrors.ErrInvalidTransactionInput
		}
		if _, err := spv.VerifyScripts(input.SourceTransaction); err != nil {
			return spverrors.Wrapf(err, "SPV script verification failed for input %d", i)
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
			return txerrors.ErrInputSourceTxIDNotFound
		}

		input.SourceTransaction = val.Tx
		if val.HadBeef {
			continue
		}

		if err := s.resolveRecursive(input.SourceTransaction.Inputs); err != nil {
			return spverrors.Wrapf(err, "Transaction %s failed to resolve source transaction for input %d", sourceTxID, idx)
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
		return nil, spverrors.Wrapf(err, "failed to convert tx query result slice into map")
	}

	return &SourceTransactionResolver{subjectTx: tx, sourceTxs: txs}, nil
}
