package engine

import (
	"context"
	"iter"
	"slices"
	"time"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

// loadFromDBTimeout - within this time transactions should be loaded from the database
const loadFromDBTimeout = 20 * time.Second

type sdkTxGetter struct {
	client *Client
}

func newSDKTxGetter(client *Client) *sdkTxGetter {
	return &sdkTxGetter{client: client}
}

func (g *sdkTxGetter) GetTransactions(ctx context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error) {
	db := g.client.Datastore().DB()

	queryIDsCtx, cancel := context.WithTimeout(ctx, loadFromDBTimeout)
	defer cancel()

	var hexes []struct {
		Hex string
	}

	idsSlice := slices.Collect(ids)
	if len(idsSlice) == 0 {
		return nil, nil
	}

	err := db.
		WithContext(queryIDsCtx).
		Model(&Transaction{}).
		Where("id IN (?)", idsSlice).
		Find(&hexes).
		Error

	if err != nil {
		return nil, spverrors.Wrapf(err, "Cannot get transactions by IDs from database")
	}

	transactions := make([]*sdk.Transaction, 0, len(hexes))
	for _, record := range hexes {
		tx, err := sdk.NewTransactionFromHex(record.Hex)
		if err != nil {
			return nil, spverrors.Wrapf(err, "Cannot parse transaction hex")
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}
