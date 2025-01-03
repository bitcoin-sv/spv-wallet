package txtracker

import (
	"context"
	"iter"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
)

// Repository is an interface for getting and saving transactions.
type Repository interface {
	MissingTransactions(ctx context.Context, txIDs iter.Seq[string]) (iter.Seq[string], error)
	SaveTXs(ctx context.Context, txRows iter.Seq[*database.TrackedTransaction]) error
}
