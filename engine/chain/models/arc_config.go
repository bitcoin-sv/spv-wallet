package chainmodels

import (
	"context"
	"iter"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
)

// TransactionsGetter is an interface for getting transactions by their IDs
type TransactionsGetter interface {
	GetTransactions(ctx context.Context, ids iter.Seq[string]) ([]*sdk.Transaction, error)
}

// ARCCallbackConfig is the configuration for spv-wallet's endpoint for ARC to callback.
type ARCCallbackConfig struct {
	URL   string
	Token string
}

// ARCConfig is the configuration for the ARC API.
type ARCConfig struct {
	URL          string
	Token        string
	DeploymentID string
	Callback     *ARCCallbackConfig
	UseJunglebus bool
	TxsGetter    TransactionsGetter
}
