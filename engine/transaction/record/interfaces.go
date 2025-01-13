package record

import (
	"context"
	"iter"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

type AddressesRepo interface {
	FindByStringAddresses(ctx context.Context, addresses iter.Seq[string]) ([]*database.Address, error)
}

type OutputsRepo interface {
	FindByOutpoints(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database.UserUtxos, []*database.TrackedOutput, error)
}

type OperationsRepo interface {
	SaveAll(ctx context.Context, opRows iter.Seq[*database.Operation]) error
}

// Broadcaster is an interface for broadcasting transactions.
type Broadcaster interface {
	Broadcast(ctx context.Context, tx *trx.Transaction) (*chainmodels.TXInfo, error)
}
