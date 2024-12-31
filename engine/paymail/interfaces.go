package paymail

import (
	"context"
	"github.com/bitcoin-sv/go-paymail/spv"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/database"
)

// Repository is an interface for the paymail repository
type Repository interface {
	GetPaymailByAlias(alias, domain string) (*database.Paymail, error)
	SaveAddress(ctx context.Context, userRow *database.User, addressRow *database.Address) error
}

type MerkleRootsVerifier interface {
	VerifyMerkleRoots(ctx context.Context, merkleRoots []*spv.MerkleRootConfirmationRequestItem) (bool, error)
}

type TxRecorder interface {
	RecordTransaction(ctx context.Context, tx *trx.Transaction, verifyScripts bool) error
}
