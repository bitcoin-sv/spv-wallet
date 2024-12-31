package txtracker

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"iter"
)

type Repository interface {
	MissingTransactions(ctx context.Context, txIDs iter.Seq[string]) (iter.Seq[string], error)
	SaveTXs(ctx context.Context, txRows iter.Seq[*database.TrackedTransaction]) error
}
