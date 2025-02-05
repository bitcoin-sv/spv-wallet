package operations

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/operations/operationsmodels"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// Repo is an interface for operations repository.
type Repo interface {
	PaginatedForUser(ctx context.Context, userID string, page filter.Page) (*models.PagedResult[operationsmodels.Operation], error)
}
