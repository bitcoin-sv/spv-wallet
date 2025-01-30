package data

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/data/datamodels"
)

// Service is the domain service for data.
type Service struct {
	dataRepo Repo
}

// NewService creates a new instance of the data service.
func NewService(dataRepo Repo) *Service {
	return &Service{
		dataRepo: dataRepo,
	}
}

// FindForUser returns the data by outpoint for a specific user.
func (s *Service) FindForUser(ctx context.Context, id string, userID string) (*datamodels.Data, error) {
	item, err := s.dataRepo.FindForUser(ctx, id, userID)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to find data for user %s", userID)
	}
	return item, nil
}
