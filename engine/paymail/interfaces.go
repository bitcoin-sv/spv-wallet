package paymail

import (
	"context"
	"iter"

	"github.com/bitcoin-sv/go-paymail/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
)

// Repository is an interface for the paymail repository
type Repository interface {
	GetPaymailByAlias(alias, domain string) (*database.Paymail, error)
	SaveAddress(ctx context.Context, userRow *database.User, addressRow *database.Address) error
}

// MerkleRootsVerifier is an interface for verifying merkle roots
type MerkleRootsVerifier interface {
	VerifyMerkleRoots(ctx context.Context, merkleRoots []*spv.MerkleRootConfirmationRequestItem) (bool, error)
}

// TxRecorder is an interface for recording transactions
type TxRecorder interface {
	RecordTransaction(ctx context.Context, tx *trx.Transaction, verifyScripts bool) error
}

// TxTracker is an interface for storing missing transactions
type TxTracker interface {
	TrackMissingTxs(ctx context.Context, transactions iter.Seq[*trx.Transaction]) error
}
