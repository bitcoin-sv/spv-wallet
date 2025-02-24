package record

import (
	"context"
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"iter"

	"github.com/bitcoin-sv/go-paymail"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses/addressesmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/beef"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// AddressesService is an interface for addresses service.
type AddressesService interface {
	FindByStringAddresses(ctx context.Context, addresses iter.Seq[string]) ([]addressesmodels.Address, error)
}

type UsersService interface {
	GetPubKey(ctx context.Context, userID string) (*primitives.PublicKey, error)
}

// OutputsRepo is an interface for outputs repository.
type OutputsRepo interface {
	FindByOutpoints(ctx context.Context, outpoints iter.Seq[bsv.Outpoint]) ([]txmodels.TrackedOutput, error)
}

// TransactionsRepo is an interface for transactions repository.
type TransactionsRepo interface {
	// FindTransactionInputSources retrieves the full ancestry of input sources for a given transaction.
	// It recursively traces input sources in batches to minimize database queries.
	FindTransactionInputSources(ctx context.Context, sourceTXIDs ...string) (beef.TxQueryResultSlice, error)
	// HasTransactionInputSources checks if all the provided input source transaction IDs exist in the database.
	// If all of them are found, the transaction data can be serialized into Raw HEX format.
	// Otherwise, serialization should be done using the BEEFHex format.
	HasTransactionInputSources(ctx context.Context, sourceTXIDs ...string) (bool, error)
}

// OperationsRepo is an interface for operations repository.
type OperationsRepo interface {
	SaveAll(ctx context.Context, opRows iter.Seq[*txmodels.NewOperation]) error
}

// Broadcaster is an interface for broadcasting transactions.
type Broadcaster interface {
	Broadcast(ctx context.Context, tx *trx.Transaction) (*chainmodels.TXInfo, error)
}

// PaymailNotifier is an interface for notifying paymail recipients about incoming transactions.
type PaymailNotifier interface {
	Notify(ctx context.Context, address string, p2pMetadata *paymail.P2PMetaData, reference string, tx *trx.Transaction) error
}
