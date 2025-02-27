package txsync

import (
	"context"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Handle(ctx context.Context, txInfo chainmodels.TXInfo) error {
	return nil
}
