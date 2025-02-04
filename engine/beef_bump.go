package engine

import (
	"context"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/samber/lo"
)

type beefPreparationContext struct {
	ctx                 context.Context
	store               TransactionGetter
	transactions        map[chainhash.Hash]*trx.Transaction
	visitedTransactions map[chainhash.Hash]bool
	queue               []*trx.Transaction
	txsForBEEF          []*trx.Transaction
}

func prepareBEEFFactors(ctx context.Context, tx *Transaction, store TransactionGetter) ([]*trx.Transaction, error) {

	processedTx, err := trx.NewTransactionFromHex(tx.Hex)
	if err != nil {
		return nil, spverrors.Wrapf(err, "cannot convert processed transaction to SDK transaction from hex (tx.ID: %s)", tx.ID)
	}

	bctx := &beefPreparationContext{
		ctx:                 ctx,
		store:               store,
		transactions:        make(map[chainhash.Hash]*trx.Transaction),
		visitedTransactions: make(map[chainhash.Hash]bool),
		queue:               []*trx.Transaction{},
		txsForBEEF:          []*trx.Transaction{processedTx},
	}

	processedSDKTx := processedTx
	txID := *processedSDKTx.TxID()
	bctx.transactions[txID] = processedSDKTx
	bctx.queue = append(bctx.queue, processedSDKTx)

	for len(bctx.queue) > 0 {
		currentTx := bctx.queue[0]
		bctx.queue = bctx.queue[1:]

		if err := bctx.processTransaction(currentTx); err != nil {
			return nil, err
		}
	}

	return bctx.txsForBEEF, nil
}

func (bctx *beefPreparationContext) processTransaction(currentTx *trx.Transaction) error {
	currentTxID := *currentTx.TxID()

	if bctx.visitedTransactions[currentTxID] {
		return nil
	}
	bctx.visitedTransactions[currentTxID] = true

	var inputTxIDs []string
	for _, input := range currentTx.Inputs {
		if input.SourceTXID != nil {
			inputTxIDs = append(inputTxIDs, input.SourceTXID.String())
		}
	}

	return bctx.processInputTransactions(inputTxIDs, currentTx)
}

func (bctx *beefPreparationContext) processInputTransactions(inputTxIDs []string, currentTx *trx.Transaction) error {
	inputTxIDs = lo.Uniq(inputTxIDs)

	inputTxs, err := bctx.store.GetTransactionsByIDs(bctx.ctx, inputTxIDs)
	if err != nil {
		return spverrors.Wrapf(err, "cannot get transactions from database")
	}

	if len(inputTxs) != len(inputTxIDs) {
		missingTxIDs := getMissingTxs(inputTxIDs, inputTxs)
		return spverrors.Newf("required transactions (%v) not found in database", missingTxIDs)
	}

	for _, inputTx := range inputTxs {
		inputSDKTx, err := convertToSDKTransaction(inputTx)
		if err != nil {
			return err
		}

		inputTxIDHash := *inputSDKTx.TxID()

		if _, exists := bctx.transactions[inputTxIDHash]; !exists {
			bctx.transactions[inputTxIDHash] = inputSDKTx
			bctx.txsForBEEF = append(bctx.txsForBEEF, inputSDKTx)
		}

		linkSourceTransaction(currentTx, inputSDKTx)
		if inputTx.BUMP.BlockHeight == 0 && len(inputTx.BUMP.Path) == 0 {
			bctx.queue = append(bctx.queue, inputSDKTx)
		}
	}

	return nil
}

func convertToSDKTransaction(inputTx *Transaction) (*trx.Transaction, error) {
	inputSDKTx, err := trx.NewTransactionFromHex(inputTx.Hex)
	if err != nil {
		return nil, spverrors.Wrapf(err, "cannot create SDK transaction from hex (tx.ID: %s)", inputTx.ID)
	}

	if inputTx.BUMP.BlockHeight != 0 {
		merklePath, err := inputTx.BUMP.toMerklePath()
		if err != nil {
			return nil, spverrors.Wrapf(err, "cannot convert BUMP to MerklePath (tx.ID: %s)", inputTx.ID)
		}
		inputSDKTx.MerklePath = merklePath
	}

	return inputSDKTx, nil
}

func linkSourceTransaction(currentTx, inputSDKTx *trx.Transaction) {
	inputTxIDHash := *inputSDKTx.TxID()
	for _, input := range currentTx.Inputs {
		if input.SourceTXID != nil && *input.SourceTXID == inputTxIDHash {
			input.SourceTransaction = inputSDKTx
		}
	}
}

func getMissingTxs(txIDs []string, foundTxs []*Transaction) []string {
	foundTxIDSet := make(map[string]struct{}, len(foundTxs))
	for _, tx := range foundTxs {
		foundTxIDSet[tx.ID] = struct{}{}
	}

	var missingTxIDs []string
	for _, txID := range txIDs {
		if _, found := foundTxIDSet[txID]; !found {
			missingTxIDs = append(missingTxIDs, txID)
		}
	}
	return missingTxIDs
}
