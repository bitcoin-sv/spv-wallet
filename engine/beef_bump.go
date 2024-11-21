package engine

import (
	"context"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

type beefPreparationContext struct {
	ctx            context.Context
	store          TransactionGetter
	transactionMap map[chainhash.Hash]*trx.Transaction
	visited        map[chainhash.Hash]bool
	queue          []*trx.Transaction
	txsForBEEF     []*trx.Transaction
}

func prepareBEEFFactors(ctx context.Context, tx *Transaction, store TransactionGetter) ([]*trx.Transaction, error) {
	txsForBEEF, err := initializeRequiredTxsCollection(tx)
	if err != nil {
		return nil, err
	}

	bctx := &beefPreparationContext{
		ctx:            ctx,
		store:          store,
		transactionMap: make(map[chainhash.Hash]*trx.Transaction),
		visited:        make(map[chainhash.Hash]bool),
		queue:          []*trx.Transaction{},
		txsForBEEF:     txsForBEEF,
	}

	processedSDKTx := txsForBEEF[0]
	txID := *processedSDKTx.TxID()
	bctx.transactionMap[txID] = processedSDKTx
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

	if bctx.visited[currentTxID] {
		return nil
	}
	bctx.visited[currentTxID] = true

	var inputTxIDs []string
	for _, input := range currentTx.Inputs {
		if input.SourceTXID != nil {
			inputTxIDs = append(inputTxIDs, input.SourceTXID.String())
		}
	}

	return bctx.retrieveAndProcessInputTransactions(inputTxIDs, currentTx)
}

func (bctx *beefPreparationContext) retrieveAndProcessInputTransactions(inputTxIDs []string, currentTx *trx.Transaction) error {
	inputTxs, err := getRequiredTransactions(bctx.ctx, inputTxIDs, bctx.store)
	if err != nil {
		return err
	}

	for _, inputTx := range inputTxs {
		inputSDKTx, err := convertToSDKTransaction(inputTx)
		if err != nil {
			return err
		}
		inputTxIDHash := *inputSDKTx.TxID()

		if _, exists := bctx.transactionMap[inputTxIDHash]; !exists {
			bctx.transactionMap[inputTxIDHash] = inputSDKTx
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
		merklePath, err := inputTx.BUMP.ToMerklePath()
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

func initializeRequiredTxsCollection(tx *Transaction) ([]*trx.Transaction, error) {
	processedSDKTx, err := trx.NewTransactionFromHex(tx.Hex)
	if err != nil {
		return nil, spverrors.Wrapf(err, "cannot convert processed transaction to SDK transaction from hex (tx.ID: %s)", tx.ID)
	}
	return []*trx.Transaction{processedSDKTx}, nil
}
