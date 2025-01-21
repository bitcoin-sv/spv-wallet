package paymail

import (
	"context"

	"github.com/bitcoin-sv/go-paymail/spv"
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/domainmodels"
)

// PaymailsRepo is an interface for paymails repository.
type PaymailsRepo interface {
	Get(ctx context.Context, alias, domain string) (*domainmodels.Paymail, error)
}

// UsersService is an interface for user service
type UsersService interface {
	AppendAddress(ctx context.Context, newAddress domainmodels.NewAddress) error
	GetPubKey(ctx context.Context, userID string) (*primitives.PublicKey, error)
}

// MerkleRootsVerifier is an interface for verifying merkle roots
type MerkleRootsVerifier interface {
	VerifyMerkleRoots(ctx context.Context, merkleRoots []*spv.MerkleRootConfirmationRequestItem) (bool, error)
}

// TxRecorder is an interface for recording transactions
type TxRecorder interface {
	RecordPaymailTransaction(ctx context.Context, tx *trx.Transaction, senderPaymail, receiverPaymail string) error
}
