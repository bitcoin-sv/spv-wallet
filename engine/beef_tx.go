package engine

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ToBeef generates BEEF Hex for transaction
func ToBeef(ctx context.Context, tx *Transaction, store TransactionGetter) (string, error) {
	if err := hydrateTransaction(ctx, tx); err != nil {
		return "", err
	}

	txsForBEEF, err := prepareBEEFFactors(ctx, tx, store)
	if err != nil {
		return "", spverrors.Wrapf(err, "prepareBUMPFactors() error")
	}

	beefBytes, err := txsForBEEF[0].BEEFHex()
	if err != nil {
		return "", spverrors.Wrapf(err, "BEEF() error")
	}

	return beefBytes, nil
}

func hydrateTransaction(ctx context.Context, tx *Transaction) error {
	if tx.draftTransaction == nil {
		dTx, err := getDraftTransactionID(
			ctx, tx.XPubID, tx.DraftID, tx.GetOptions(false)...,
		)

		if err != nil || dTx == nil {
			return spverrors.Wrapf(err, "retrieve DraftTransaction failed")
		}

		tx.draftTransaction = dTx
	}

	return nil
}
