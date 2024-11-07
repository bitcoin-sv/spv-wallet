package engine

import (
	"context"
	"encoding/hex"

	chainhash "github.com/bitcoin-sv/go-sdk/chainhash"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

const maxBeefVer = uint32(0xFFFF) // value from BRC-62

type beefTx struct {
	version      uint32
	bumps        BUMPs
	transactions []*trx.Transaction
}

// ToBeef generates BEEF Hex for transaction
func ToBeef(ctx context.Context, tx *Transaction, store TransactionGetter) (string, error) {
	if err := hydrateTransaction(ctx, tx); err != nil {
		return "", err
	}

	bumpBtFactors, bumpFactors, err := prepareBEEFFactors(ctx, tx, store)
	if err != nil {
		return "", spverrors.Wrapf(err, "prepareBUMPFactors() error")
	}

	bumps, err := calculateMergedBUMP(bumpFactors)
	if err != nil {
		return "", err
	}
	sortedTxs := kahnTopologicalSortTransactions(bumpBtFactors)
	beefHex, err := toBeefHex(bumps, sortedTxs)
	if err != nil {
		return "", spverrors.Wrapf(err, "ToBeef() error")
	}

	return beefHex, nil
}

func toBeefHex(bumps BUMPs, parentTxs []*trx.Transaction) (string, error) {
	beef, err := newBeefTx(1, bumps, parentTxs)
	if err != nil {
		return "", spverrors.Wrapf(err, "ToBeefHex() error")
	}

	beefBytes, err := beef.toBeefBytes()
	if err != nil {
		return "", spverrors.Wrapf(err, "ToBeefHex() error")
	}

	return hex.EncodeToString(beefBytes), nil
}

func newBeefTx(version uint32, bumps BUMPs, parentTxs []*trx.Transaction) (*beefTx, error) {
	if version > maxBeefVer {
		return nil, spverrors.Newf("version above 0x%X", maxBeefVer)
	}

	if err := validateBumps(bumps); err != nil {
		return nil, err
	}

	beef := &beefTx{
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
