package record

import (
	"context"
	"iter"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	database2 "github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// AddressesRepo is an interface for addresses repository.
type AddressesRepo interface {
	FindByStringAddresses(ctx context.Context, addresses iter.Seq[string]) ([]*database2.Address, error)
}

// OutputsRepo is an interface for outputs repository.
type OutputsRepo interface {
	FindByOutpoints(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]*database2.TrackedOutput, error)
}

// OperationsRepo is an interface for operations repository.
type OperationsRepo interface {
	SaveAll(ctx context.Context, opRows iter.Seq[*database2.Operation]) error
}

// Broadcaster is an interface for broadcasting transactions.
type Broadcaster interface {
	Broadcast(ctx context.Context, tx *trx.Transaction) (*chainmodels.TXInfo, error)
}
