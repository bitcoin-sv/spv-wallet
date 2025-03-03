package txmodels

import (
	"time"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

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

func (tt *TrackedTransaction) Mined(blockHash string, blockHeight int64, bump *trx.MerklePath) error {
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
	tt.BlockHeight = &blockHeight
	tt.TxStatus = TxStatusMined

	return nil
}
