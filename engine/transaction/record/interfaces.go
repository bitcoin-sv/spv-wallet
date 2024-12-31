package record

import (
	"context"
	"iter"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// Repository is an interface for saving transactions and outputs to the database.
type Repository interface {
	SaveTX(ctx context.Context, txRow *database.TrackedTransaction) error
	GetOutputs(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.Output, error)
	CheckAddress(ctx context.Context, address string) (bool, error)
}

// Broadcaster is an interface for broadcasting transactions.
type Broadcaster interface {
	Broadcast(ctx context.Context, tx *trx.Transaction) (*chainmodels.TXInfo, error)
}
