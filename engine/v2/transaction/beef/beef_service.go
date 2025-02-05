package beef

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
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

// TxRepository defines an interface for querying transaction input sources.
// It provides a method to retrieve transaction details for a given set of source transaction IDs.
type TxRepository interface {
	// QueryInputSources retrieves transaction query results for the provided source transaction IDs.
	// Returns a slice of TxQueryResult and an error if the query fails.
	QueryInputSources(ctx context.Context, sourceTXIDs ...string) (TxQueryResultSlice, error)
}

// Service provides transaction processing functionalities, including preparing BEEF-encoded transactions.
type Service struct {
	repository TxRepository // Repository used to query transaction input sources.
}

// extractSourceTXIDs extracts source transaction IDs from the inputs of the given transaction.
// It returns a slice of transaction IDs or an error if the transaction has no inputs.
func (s *Service) extractSourceTXIDs(tx *sdk.Transaction) ([]string, error) {
	if tx.InputCount() == 0 {
		return nil, txerrors.ErrZeroInputCount
	}

	sourceTXIDs := make([]string, len(tx.Inputs))
	for _, in := range tx.Inputs {
		sourceTXIDs = append(sourceTXIDs, in.SourceTXID.String())
	}
	return sourceTXIDs, nil
}

// PrepareBEEF constructs a BEEF-encoded transaction representation for the given transaction.
// It resolves source transactions for all inputs before encoding the transaction in BEEF format.
// Returns the BEEF hex string or an error if resolution or encoding fails.
func (s *Service) PrepareBEEF(ctx context.Context, tx *sdk.Transaction) (string, error) {
	if tx == nil {
		return "", txerrors.ErrNilSubjectTx
	}

	txID := tx.TxID().String()
	sourceTxIDs, err := s.extractSourceTXIDs(tx)
	if err != nil {
		return "", spverrors.Wrapf(err, "failed to extract source transaction IDs for transaction %s", txID)
	}

	txQueryResult, err := s.repository.QueryInputSources(ctx, sourceTxIDs...)
	if err != nil {
		return "", spverrors.Wrapf(err, "database query failed while retrieving input data for transaction %s", txID)
	}

	resolver, err := NewSourceTransactionResolver(tx, txQueryResult)
	if err != nil {
		return "", spverrors.Wrapf(err, "failed to initialize source transaction resolver for transaction %s", txID)
	}

	err = resolver.Resolve()
	if err != nil {
		return "", spverrors.Wrapf(err, "failed to resolve source transactions for transaction %s", txID)
	}

	hex, err := tx.BEEFHex()
	if err != nil {
		return "", spverrors.Wrapf(err, "failed to generate BEEF hex encoding for transaction %s", txID)
	}

	return hex, nil
}

// NewService creates a new Service instance with the provided TxRepository.
// Panics if the repository is nil.
func NewService(r TxRepository) *Service {
	if r == nil {
		panic("transactions repository must be non-nil value")
	}
	return &Service{repository: r}
}
