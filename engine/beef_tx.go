package engine

import (
	"context"
	"encoding/hex"

	"github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

const maxBeefVer = uint32(0xFFFF) // value from BRC-62

type beefTxSDK struct {
	version      uint32
	bumps        BUMPs
	transactions []*trx.Transaction
}

// ToBeef generates BEEF Hex for transaction
func ToBeef(ctx context.Context, tx *Transaction, store TransactionGetter) (string, error) {
	if err := hydrateTransaction(ctx, tx); err != nil {
		return "", err
	}

	//bumpBtFactors, bumpFactors, err := prepareBEEFFactors(ctx, tx, store)
	bumpBtFactors, _, err := prepareBEEFFactors(ctx, tx, store)
	if err != nil {
		return "", spverrors.Wrapf(err, "prepareBUMPFactors() error")
	}

	bumpsSDK, err := calculateMergedBUMPSDK(bumpBtFactors)
	if err != nil {
		return "", spverrors.Wrapf(err, "calculateMergedBUMPSDK() error")
	}

	beefHex, err := toBeefHexSDK(bumpsSDK, bumpBtFactors)
	if err != nil {
		return "", spverrors.Wrapf(err, "ToBeef() error")
	}

	return beefHex, nil
}

func toBeefHexSDK(bumps []*trx.MerklePath, parentTxs []*trx.Transaction) (string, error) {
	beef, err := newBeefTxSDK(1, bumps, parentTxs)
	if err != nil {
		return "", spverrors.Wrapf(err, "ToBeefHex() error")
	}

	populateSourceTransactions(parentTxs)

	beefBytes, err := beef.transactions[0].BEEF()
	if err != nil {
		return "", spverrors.Wrapf(err, "ToBeefHex() error")
	}

	return hex.EncodeToString(beefBytes), nil
}

func newBeefTxSDK(version uint32, bumps []*trx.MerklePath, parentTxs []*trx.Transaction) (*beefTxSDK, error) {

	if version > maxBeefVer {
		return nil, spverrors.Newf("version above 0x%X", maxBeefVer)
	}

	if err := validatMerklePaths(bumps); err != nil {
		return nil, err
	}

	beef := &beefTxSDK{
		version:      version,
		bumps:        bumps,
		transactions: parentTxs,
	}

	return beef, nil
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

func populateSourceTransactions(transactions []*trx.Transaction) {
	transactionMap := make(map[chainhash.Hash]*trx.Transaction)
	for _, tx := range transactions {
		txID := *tx.TxID()
		transactionMap[txID] = tx
	}

	visited := make(map[chainhash.Hash]bool)
	queue := make([]*trx.Transaction, 0, len(transactions))

	for _, tx := range transactions {
		txID := *tx.TxID()
		if visited[txID] {
			continue
		}
		queue = append(queue, tx)

		for len(queue) > 0 {
			currentTx := queue[0]
			queue = queue[1:]
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
						queue = append(queue, sourceTransaction)
					}
				}
			}
		}
	}
}
