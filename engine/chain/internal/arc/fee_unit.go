package arc

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// GetFeeUnit returns the current fee unit from the ARC policy.
func (s *Service) GetFeeUnit(ctx context.Context) (*bsv.FeeUnit, error) {
	policy, err := s.getPolicy(ctx)
	if err != nil {
		return nil, err
	}

	return &bsv.FeeUnit{
		Satoshis: policy.Content.MiningFee.Satoshis,
		Bytes:    policy.Content.MiningFee.Bytes,
	}, nil
}
