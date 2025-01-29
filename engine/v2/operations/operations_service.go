package operations

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/operations/operationsmodels"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) PaginatedForUser(ctx context.Context, userID string, page filter.Page) (*models.PagedResult[operationsmodels.Operation], error) {
	entities, err := s.repo.PaginatedForUser(ctx, userID, page)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to get operations for user")
	}

	return entities, nil
}
