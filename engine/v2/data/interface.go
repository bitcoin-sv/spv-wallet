package data

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/data/datamodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// Repo is the interface that wraps the basic operations with data.
type Repo interface {
	FindForUser(ctx context.Context, outpoint bsv.Outpoint, userID string) (*datamodels.Data, error)
}
