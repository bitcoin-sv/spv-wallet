package engine

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/libsv/go-bt/v2"
)

const maxBeefVer = uint32(0xFFFF) // value from BRC-62

type beefTx struct {
	version      uint32
	bumps        BUMPs
	transactions []*bt.Tx
}

// ToBeef generates BEEF Hex for transaction
func ToBeef(ctx context.Context, tx *Transaction, store TransactionGetter) (string, error) {
	if err := hydrateTransaction(ctx, tx); err != nil {
		return "", err
	}

	bumpBtFactors, bumpFactors, err := prepareBEEFFactors(ctx, tx, store)
	if err != nil {
		return "", fmt.Errorf("prepareBUMPFactors() error: %w", err)
	}

	bumps, err := calculateMergedBUMP(bumpFactors)
	if err != nil {
		return "", err
	}
	sortedTxs := kahnTopologicalSortTransactions(bumpBtFactors)
	beefHex, err := toBeefHex(bumps, sortedTxs)
	if err != nil {
		return "", fmt.Errorf("ToBeef() error: %w", err)
	}

	return beefHex, nil
}

func toBeefHex(bumps BUMPs, parentTxs []*bt.Tx) (string, error) {
	beef, err := newBeefTx(1, bumps, parentTxs)
	if err != nil {
		return "", fmt.Errorf("ToBeefHex() error: %w", err)
	}

	beefBytes, err := beef.toBeefBytes()
	if err != nil {
		return "", fmt.Errorf("ToBeefHex() error: %w", err)
	}

	return hex.EncodeToString(beefBytes), nil
}

func newBeefTx(version uint32, bumps BUMPs, parentTxs []*bt.Tx) (*beefTx, error) {
	if version > maxBeefVer {
		return nil, fmt.Errorf("version above 0x%X", maxBeefVer)
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
			return fmt.Errorf("retrieve DraftTransaction failed: %w", err)
		}

		tx.draftTransaction = dTx
	}

	return nil
}
