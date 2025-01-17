package paymail

import (
	"context"

	"github.com/bitcoin-sv/go-paymail/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	paymailmodels "github.com/bitcoin-sv/spv-wallet/engine/paymail/models"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// PaymailsRepo is an interface for paymails repository.
type PaymailsRepo interface {
	Get(ctx context.Context, alias, domain string) (*paymailmodels.Paymail, error)
}

// UsersService is an interface for user service
type UsersService interface {
	AppendAddress(ctx context.Context, userID string, address string, customInstructions bsv.CustomInstructions) error
}

// MerkleRootsVerifier is an interface for verifying merkle roots
type MerkleRootsVerifier interface {
	VerifyMerkleRoots(ctx context.Context, merkleRoots []*spv.MerkleRootConfirmationRequestItem) (bool, error)
}

// TxRecorder is an interface for recording transactions
type TxRecorder interface {
	RecordPaymailTransaction(ctx context.Context, tx *trx.Transaction, senderPaymail, receiverPaymail string) error
}
