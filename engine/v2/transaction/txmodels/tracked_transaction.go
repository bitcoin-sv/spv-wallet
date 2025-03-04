package txmodels

import (
	"github.com/samber/lo"
	"time"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// TrackedTransaction represents a transaction that is being tracked by the wallet.
type TrackedTransaction struct {
	ID       string
	TxStatus TxStatus

	CreatedAt time.Time
	UpdatedAt time.Time

	BlockHeight *int64
	BlockHash   *string

	BeefHex *string
	RawHex  *string
}

// TX returns the transaction object for the tracked transaction based on the hex representation (either BEEF or raw).
func (tt *TrackedTransaction) TX() (*trx.Transaction, error) {
	var tx *trx.Transaction
	var err error
	if tt.BeefHex != nil {
		tx, err = trx.NewTransactionFromBEEFHex(*tt.BeefHex)
	} else if tt.RawHex != nil {
		tx, err = trx.NewTransactionFromHex(*tt.RawHex)
	} else {
		return nil, spverrors.Newf("tracked transaction %s has no transaction hex", tt.ID)
	}

	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to parse transaction hex for transaction %s", tt.ID)
	}

	return tx, nil
}

// Mined marks the transaction as mined with the given block hash and height, and the given bump.
func (tt *TrackedTransaction) Mined(blockHash string, bump *trx.MerklePath) error {
	tx, err := tt.TX()
	if err != nil {
		return err
	}

	tx.MerklePath = bump

	beefHex, err := tx.BEEFHex()
	if err != nil {
		return spverrors.Wrapf(err, "failed to get BEEF hex for transaction %s", tt.ID)
	}

	tt.BeefHex = &beefHex
	tt.RawHex = nil
	tt.BlockHash = &blockHash
	tt.BlockHeight = lo.ToPtr(int64(bump.BlockHeight))
	tt.TxStatus = TxStatusMined

	return nil
}
