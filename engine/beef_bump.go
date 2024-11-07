package engine

import (
	"context"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

func prepareBEEFFactors(ctx context.Context, tx *Transaction, store TransactionGetter) ([]*trx.Transaction, []*Transaction, error) {
	sdkTxsNeededForBUMP, txsNeededForBUMP, err := initializeRequiredTxsCollection(tx)
	if err != nil {
		return nil, nil, err
	}

	txIDs := make([]string, 0, len(tx.draftTransaction.Configuration.Inputs))
	for _, input := range tx.draftTransaction.Configuration.Inputs {
		txIDs = append(txIDs, input.UtxoPointer.TransactionID)
	}

	inputTxs, err := getRequiredTransactions(ctx, txIDs, store)
	if err != nil {
		return nil, nil, err
	}

	for _, inputTx := range inputTxs {
		inputSDKTx, err := trx.NewTransactionFromHex(inputTx.Hex)
		if err != nil {
			return nil, nil, spverrors.Wrapf(err, "cannot convert to SDK Tx from hex (tx.ID: %s)", inputTx.ID)
		}

		txsNeededForBUMP = append(txsNeededForBUMP, inputTx)
		sdkTxsNeededForBUMP = append(sdkTxsNeededForBUMP, inputSDKTx)

		if inputTx.BUMP.BlockHeight == 0 && len(inputTx.BUMP.Path) == 0 {
			parentSDKTransactions, parentTransactions, err := checkParentTransactions(ctx, store, inputSDKTx)
			if err != nil {
				return nil, nil, err
			}

			txsNeededForBUMP = append(txsNeededForBUMP, parentTransactions...)
			sdkTxsNeededForBUMP = append(sdkTxsNeededForBUMP, parentSDKTransactions...)
		}
	}

	return sdkTxsNeededForBUMP, txsNeededForBUMP, nil
}

func checkParentTransactions(ctx context.Context, store TransactionGetter, sdkTx *trx.Transaction) ([]*trx.Transaction, []*Transaction, error) {
	parentTxIDs := make([]string, 0, len(sdkTx.Inputs))
	for _, txIn := range sdkTx.Inputs {
		parentTxIDs = append(parentTxIDs, txIn.SourceTXID.String())
	}

	parentTxs, err := getRequiredTransactions(ctx, parentTxIDs, store)
	if err != nil {
		return nil, nil, err
	}

	validTxs := make([]*Transaction, 0, len(parentTxs))
	validSDKTxs := make([]*trx.Transaction, 0, len(parentTxs))
	for _, parentTx := range parentTxs {
		parentSDKTx, err := trx.NewTransactionFromHex(parentTx.Hex)
		if err != nil {
			return nil, nil, spverrors.Wrapf(err, "cannot convert to SDK Tx from hex (tx.ID: %s)", parentTx.ID)
		}
		validTxs = append(validTxs, parentTx)
		validSDKTxs = append(validSDKTxs, parentSDKTx)

		if parentTx.BUMP.BlockHeight == 0 && len(parentTx.BUMP.Path) == 0 {
			parentValidSDKTxs, parentValidTxs, err := checkParentTransactions(ctx, store, parentSDKTx)
			if err != nil {
				return nil, nil, err
			}
			validTxs = append(validTxs, parentValidTxs...)
			validSDKTxs = append(validSDKTxs, parentValidSDKTxs...)
		}
	}

	return validSDKTxs, validTxs, nil
}

func getRequiredTransactions(ctx context.Context, txIDs []string, store TransactionGetter) ([]*Transaction, error) {
	txs, err := store.GetTransactionsByIDs(ctx, txIDs)
	if err != nil {
		return nil, spverrors.Wrapf(err, "cannot get transactions from database")
	}

	if len(txs) != len(txIDs) {
		missingTxIDs := getMissingTxs(txIDs, txs)
		return nil, spverrors.Newf("required transactions (%v) not found in database", missingTxIDs)
	}

	return txs, nil
}

func getMissingTxs(txIDs []string, foundTxs []*Transaction) []string {
	foundTxIDs := make(map[string]bool)
	for _, tx := range foundTxs {
		foundTxIDs[tx.ID] = true
	}

	var missingTxIDs []string
	for _, txID := range txIDs {
		if !foundTxIDs[txID] {
			missingTxIDs = append(missingTxIDs, txID)
		}
	}
	return missingTxIDs
}

func initializeRequiredTxsCollection(tx *Transaction) ([]*trx.Transaction, []*Transaction, error) {
	var sdkTxsNeededForBUMP []*trx.Transaction
	var txsNeededForBUMP []*Transaction

	processedSDKTx, err := trx.NewTransactionFromHex(tx.Hex)
	if err != nil {
		return nil, nil, spverrors.Wrapf(err, "cannot convert processed tx to SDK Tx from hex (tx.ID: %s)", tx.ID)
	}

	sdkTxsNeededForBUMP = append(sdkTxsNeededForBUMP, processedSDKTx)
	txsNeededForBUMP = append(txsNeededForBUMP, tx)

	return sdkTxsNeededForBUMP, txsNeededForBUMP, nil
}
