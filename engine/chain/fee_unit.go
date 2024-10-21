package chain

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// GetFeeUnit returns the current fee unit from the ARC policy.
func (s *chainService) GetFeeUnit(ctx context.Context) (*bsv.FeeUnit, error) {
	policy, err := s.arcService.GetPolicy(ctx)
	if err != nil {
		return nil, chainerrors.ErrGetFeeUnit.Wrap(err)
	}

	return &bsv.FeeUnit{
		Satoshis: policy.Content.MiningFee.Satoshis,
		Bytes:    policy.Content.MiningFee.Bytes,
	}, nil
}
