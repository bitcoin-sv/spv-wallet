package chainstate

import (
	"context"
	"net/http"
	"time"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/rs/zerolog"
)

// HTTPInterface is the HTTP client interface
type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// ChainService is the chain related methods
type ChainService interface {
	SupportedBroadcastFormats() HexFormatFlag
	Broadcast(ctx context.Context, id, txHex string, format HexFormatFlag, timeout time.Duration) *BroadcastFailure
	QueryTransaction(
		ctx context.Context, id string, requiredIn RequiredIn, timeout time.Duration,
	) (*TransactionInfo, error)
	QueryTransactionFastest(
		ctx context.Context, id string, requiredIn RequiredIn, timeout time.Duration,
	) (*TransactionInfo, error)
}

// ProviderServices is the chainstate providers interface
type ProviderServices interface {
	BroadcastClient() broadcast.Client
}

// HeaderService is header services interface
type HeaderService interface {
	VerifyMerkleRoots(ctx context.Context, merkleRoots []MerkleRootConfirmationRequestItem) error
}

// ClientInterface is the chainstate client interface
type ClientInterface interface {
	ChainService
	ProviderServices
	HeaderService
	Close(ctx context.Context)
	HTTPClient() HTTPInterface
	IsNewRelicEnabled() bool
	Network() Network
	QueryTimeout() time.Duration
	FeeUnit() *utils.FeeUnit
	Logger() *zerolog.Logger
}
