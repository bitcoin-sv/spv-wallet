package engine

import (
	"context"

	chainhash "github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// ToBeef generates BEEF Hex for transaction
func ToBeef(ctx context.Context, tx *Transaction, store TransactionGetter) (string, error) {
	if err := hydrateTransaction(ctx, tx); err != nil {
		return "", err
	}

	bumpBtFactors, bumpFactors, err := prepareBEEFFactors(ctx, tx, store)
	if err != nil {
		return "", spverrors.Wrapf(err, "prepareBUMPFactors() error")
	}

	err = setMerklePathsFromBUMPs(bumpBtFactors, bumpFactors)
	if err != nil {
		return "", spverrors.Wrapf(err, "SetMerklePathsFromBUMPs() error")
	}
	populateSourceTransactions(bumpBtFactors)

	trxHex, err := bumpBtFactors[0].BEEFHex()
	if err != nil {
		return "", spverrors.Wrapf(err, "BEEFHex() error")
	}

	return trxHex, nil
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

func setMerklePathsFromBUMPs(bumpBtFactors []*trx.Transaction, bumpFactors []*Transaction) error {
	for i := range bumpBtFactors {
		if bumpFactors[i].BUMP.BlockHeight != 0 {
			merklePath, err := buildMerklePathFromBUMP(&bumpFactors[i].BUMP)
			if err != nil {
				return err
			}
			bumpBtFactors[i].MerklePath = merklePath
		}
	}
	return nil
}

func populateSourceTransactions(transactions []*trx.Transaction) {
	transactionMap := make(map[chainhash.Hash]*trx.Transaction)
	for _, tx := range transactions {
		txID := *tx.TxID()
		transactionMap[txID] = tx
	}

	visited := make(map[chainhash.Hash]bool)
	stack := make([]*trx.Transaction, 0, len(transactions))

	for _, tx := range transactions {
		txID := *tx.TxID()
		if visited[txID] {
			continue
		}
		stack = append(stack, tx)

		for len(stack) > 0 {
			currentTx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			currentTxID := *currentTx.TxID()

			if visited[currentTxID] {
				continue
			}
			visited[currentTxID] = true

			for _, input := range currentTx.Inputs {
				if input.SourceTXID != nil {
					sourceTxID := *input.SourceTXID
					if sourceTransaction, exists := transactionMap[sourceTxID]; exists {
						input.SourceTransaction = sourceTransaction
						stack = append(stack, sourceTransaction)
					}
				}
			}
		}
	}
}
