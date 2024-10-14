package engine

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
)

// ToEfHex generates Extended Format hex of transaction
func ToEfHex(ctx context.Context, tx *Transaction, store TransactionGetter) (efHex string, ok bool) {
	sdkTx := tx.parsedTx

	if sdkTx == nil {
		var err error
		sdkTx, err = trx.NewTransactionFromHex(tx.Hex)
		if err != nil {
			return "", false
		}
	}

	needToHydrate := false
	for _, input := range sdkTx.Inputs {
		if input.SourceTXID == nil {
			needToHydrate = true
			break
		}
	}

	if needToHydrate {
		if ok := hydrate(ctx, sdkTx, store); !ok {
			return "", false
		}
	}

	ef, err := sdkTx.EFHex()
	if err != nil {
		return "", false
	}

	return ef, true
}

func hydrate(ctx context.Context, tx *trx.Transaction, store TransactionGetter) (ok bool) {
	txToGet := make([]string, 0, len(tx.Inputs))

	for _, input := range tx.Inputs {
		txToGet = append(txToGet, input.SourceTXID.String())
	}

	parentTxs, err := store.GetTransactionsByIDs(ctx, txToGet)
	if err != nil {
		return false
	}
	if len(parentTxs) != len(tx.Inputs) {
		return false
	}

	for _, input := range tx.Inputs {
		prevTxID := input.SourceTXID.String()
		pTx := find(parentTxs, func(tx *Transaction) bool { return tx.ID == prevTxID })

		pbtTx, err := trx.NewTransactionFromHex((*pTx).Hex)
		if err != nil {
			return false
		}

		o := pbtTx.Outputs[input.SourceTxOutIndex]
		input.SetSourceTxOutput(&trx.TransactionOutput{
			Satoshis:      o.Satoshis,
			LockingScript: o.LockingScript,
		})
	}

	return true
}
