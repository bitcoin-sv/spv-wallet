package engine

import (
	"context"
	"slices"
	"time"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"iter"
)

type sdkTxGetter struct {
	client *Client
}

func newSDKTxGetter(client *Client) *sdkTxGetter {
	return &sdkTxGetter{client: client}
}

func (g *sdkTxGetter) GetTransactions(ctx context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error) {
	db := g.client.Datastore().DB()

	queryIDsCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var hexes []struct {
		Hex string
	}

	err := db.
		WithContext(queryIDsCtx).
		Model(&Transaction{}).
		Where("id IN (?)", slices.Collect(ids)).
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
