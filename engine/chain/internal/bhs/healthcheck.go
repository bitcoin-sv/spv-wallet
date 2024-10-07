package bhs

import (
	"context"
	"time"

	"github.com/bitcoin-sv/go-paymail/spv"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

var pastMerkleRoot = []*spv.MerkleRootConfirmationRequestItem{{
	// This is a known merkle root from the blockchain
	// timestamp: 2024-10-04 10:55:37 UTC
	MerkleRoot:  "f67ae53720205a55f4e99c632debabb68b6df0dc0f68affd200a076aee6e80e6",
	BlockHeight: 864921,
}}

// HealthcheckBHS checks the health of the Block Headers Service (BHS).
// It verifies the BHS is reachable and ready to verify merkle roots by providing a known merkle root.
func (s *Service) HealthcheckBHS(ctx context.Context) error {
	timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	confirmations, err := s.VerifyMerkleRoots(timedCtx, pastMerkleRoot)
	if err != nil {
		return chainerrors.ErrBHSUnhealthy.Wrap(err)
	}
	if confirmations.ConfirmationState != chainmodels.MRConfirmed {
		return chainerrors.ErrBHSUnhealthy.Wrap(spverrors.Newf("BHS is reachable but it's not ready (yet) to verify merkle roots"))
	}
	return nil
}
