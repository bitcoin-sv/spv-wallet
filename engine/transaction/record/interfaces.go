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
	SaveOperations(ctx context.Context, opRows iter.Seq[*database.Operation]) error
	GetOutputs(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.UserUtxos, []*database.TrackedOutput, error)
	CheckAddress(ctx context.Context, address string) (*database.Address, error)
}

// Broadcaster is an interface for broadcasting transactions.
type Broadcaster interface {
	Broadcast(ctx context.Context, tx *trx.Transaction) (*chainmodels.TXInfo, error)
}
