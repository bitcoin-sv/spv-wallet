package chainstate

import (
	"context"
	"net/http"
	"time"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// HTTPInterface is the HTTP client interface
type HTTPInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// ChainService is the chain related methods
type ChainService interface {
	SupportedBroadcastFormats() HexFormatFlag
	Broadcast(ctx context.Context, id, txHex string, format HexFormatFlag, timeout time.Duration) *BroadcastResult
	QueryTransaction(
		ctx context.Context, id string, requiredIn RequiredIn, timeout time.Duration,
	) (*TransactionInfo, error)
}

// ProviderServices is the chainstate providers interface
type ProviderServices interface {
	BroadcastClient() broadcast.Client
}

// ClientInterface is the chainstate client interface
type ClientInterface interface {
	ChainService
	ProviderServices
	Debug(on bool)
	DebugLog(text string)
	HTTPClient() HTTPInterface
	IsDebug() bool
	Network() Network
	QueryTimeout() time.Duration
	FeeUnit() *bsv.FeeUnit
}
