package record

import (
	"context"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"iter"
)

type DomainTX struct {
	ID string `gorm:"id"`
}

type Repository interface {
	SaveTX(ctx context.Context, txTable *database.Transaction, outputs []database.Output, data []database.Data) error
	GetOutputs(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]database.Output, error)
}

type Broadcaster interface {
	Broadcast(ctx context.Context, tx *trx.Transaction) error
}
