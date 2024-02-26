package engine

import (
	"context"
	"encoding/hex"

	"github.com/libsv/go-bt/v2"
)

// ToEfHex generates Extended Format hex of transaction
func ToEfHex(ctx context.Context, tx *Transaction, store TransactionGetter) (efHex string, ok bool) {
	btTx := tx.parsedTx

	if btTx == nil {
		var err error
		btTx, err = bt.NewTxFromString(tx.Hex)
		if err != nil {
			return "", false
		}
	}

	needToHydrate := false
	for _, input := range btTx.Inputs {
		if input.PreviousTxScript == nil {
			needToHydrate = true
			break
		}
	}

	if needToHydrate {
		if ok := hydrate(ctx, btTx, store); !ok {
			return "", false
		}
	}

	return hex.EncodeToString(btTx.ExtendedBytes()), true
}

func hydrate(ctx context.Context, tx *bt.Tx, store TransactionGetter) (ok bool) {
	txToGet := make([]string, 0, len(tx.Inputs))

	for _, input := range tx.Inputs {
		txToGet = append(txToGet, input.PreviousTxIDStr())
	}

	parentTxs, err := store.GetTransactionsByIDs(ctx, txToGet)
	if err != nil {
		return false
	}
	if len(parentTxs) != len(tx.Inputs) {
		return false
	}

	for _, input := range tx.Inputs {
		prevTxID := input.PreviousTxIDStr()
		pTx := find(parentTxs, func(tx *Transaction) bool { return tx.ID == prevTxID })

		pbtTx, err := bt.NewTxFromString((*pTx).Hex)
		if err != nil {
			return false
		}

		o := pbtTx.Outputs[input.PreviousTxOutIndex]
		input.PreviousTxSatoshis = o.Satoshis
		input.PreviousTxScript = o.LockingScript
	}

	return true
}
