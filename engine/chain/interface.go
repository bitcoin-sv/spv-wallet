package chain

import (
	"context"
	"github.com/bitcoin-sv/go-paymail/spv"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
)

// ARCService for querying ARC server.
type ARCService interface {
	QueryTransaction(ctx context.Context, txID string) (*chainmodels.TXInfo, error)
	GetPolicy(ctx context.Context) (*chainmodels.Policy, error)
}

// BHSService for querying BHS server.
type BHSService interface {
	VerifyMerkleRoots(ctx context.Context, merkleRoots []*spv.MerkleRootConfirmationRequestItem) (*chainmodels.MerkleRootsConfirmations, error)
	HealthcheckBHS(ctx context.Context) error
}

// Service related to the chain.
type Service interface {
	ARCService
	BHSService
}
