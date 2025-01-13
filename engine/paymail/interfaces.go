package paymail

import (
	"context"

	"github.com/bitcoin-sv/go-paymail/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
)

type PaymailsRepo interface {
	Get(ctx context.Context, alias, domain string) (*database.Paymail, error)
}

type UsersRepo interface {
	AppendAddress(ctx context.Context, userRow *database.User, addressRow *database.Address) error
}

// MerkleRootsVerifier is an interface for verifying merkle roots
type MerkleRootsVerifier interface {
	VerifyMerkleRoots(ctx context.Context, merkleRoots []*spv.MerkleRootConfirmationRequestItem) (bool, error)
}

// TxRecorder is an interface for recording transactions
type TxRecorder interface {
	RecordTransaction(ctx context.Context, tx *trx.Transaction, verifyScripts bool) error
}
