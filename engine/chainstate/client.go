package chainstate

import (
	"context"
	"time"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/engine/logging"
	"github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/rs/zerolog"
)

type (

	// Client is the client (configuration)
	Client struct {
		options *clientOptions
	}

	// clientOptions holds all the configuration for the client
	clientOptions struct {
		config    *syncConfig      // Configuration for broadcasting and other chain-state actions
		debug     bool             // For extra logs and additional debug information
		logger    *zerolog.Logger  // Logger interface
		metrics   *metrics.Metrics // For collecting metrics (if enabled)
		userAgent string           // Custom user agent for outgoing HTTP Requests
	}

	// syncConfig holds all the configuration about the different sync processes
	syncConfig struct {
		callbackURL           string           // Broadcast callback URL
		callbackToken         string           // Broadcast callback access token
		excludedProviders     []string         // List of provider names
		httpClient            HTTPInterface    // Custom HTTP client (for example WOC)
		broadcastClientConfig *broadcastConfig // Broadcast client configuration
		network               Network          // Current network (mainnet, testnet, stn)
		queryTimeout          time.Duration    // Timeout for transaction query
		broadcastClient       broadcast.Client // Broadcast client
		feeUnit               *bsv.FeeUnit     // The lowest fees among all miners
		feeQuotes             bool             // If set, feeUnit will be updated with fee quotes from miner's
	}

	broadcastConfig struct {
		ArcAPIs []string
	}
)

// NewClient creates a new client for all on-chain functionality
//
// If no options are given, it will use the defaultClientOptions()
func NewClient(ctx context.Context, opts ...ClientOps) (ClientInterface, error) {
	// Create a new client with defaults
	client := &Client{options: defaultClientOptions()}

	// Overwrite defaults with any set by user
	for _, opt := range opts {
		opt(client.options)
	}

	// Set logger if not set
	if client.options.logger == nil {
		client.options.logger = logging.GetDefaultLogger()
	}

	if err := client.initActiveProvider(ctx); err != nil {
		return nil, err
	}

	if err := client.checkFeeUnit(); err != nil {
		return nil, err
	}

	// Return the client
	return client, nil
}

// Debug will set the debug flag
func (c *Client) Debug(on bool) {
	c.options.debug = on
}

// DebugLog will display verbose logs
func (c *Client) DebugLog(text string) {
	c.options.logger.Debug().Msg(text)
}

// IsDebug will return if debugging is enabled
func (c *Client) IsDebug() bool {
	return c.options.debug
}

// HTTPClient will return the HTTP client
func (c *Client) HTTPClient() HTTPInterface {
	return c.options.config.httpClient
}

// Network will return the current network
func (c *Client) Network() Network {
	return c.options.config.network
}

// BroadcastClient will return the BroadcastClient client
func (c *Client) BroadcastClient() broadcast.Client {
	return c.options.config.broadcastClient
}

// QueryTimeout will return the query timeout
func (c *Client) QueryTimeout() time.Duration {
	return c.options.config.queryTimeout
}

// FeeUnit will return feeUnit
func (c *Client) FeeUnit() *bsv.FeeUnit {
	return c.options.config.feeUnit
}

func (c *Client) isFeeQuotesEnabled() bool {
	return c.options.config.feeQuotes
}

func (c *Client) initActiveProvider(ctx context.Context) error {
	return c.broadcastClientInit(ctx)
}

func (c *Client) checkFeeUnit() error {
	feeUnit := c.options.config.feeUnit
	switch {
	case feeUnit == nil:
		return spverrors.Newf("no fee unit found")
	case !feeUnit.IsValid():
		return spverrors.Newf("invalid fee unit found: %s", feeUnit)
	case feeUnit.IsZero():
		c.options.logger.Warn().Msg("fee unit suggests no fees (free)")
	default:
		var feeUnitSource string
		if c.isFeeQuotesEnabled() {
			feeUnitSource = "fee quotes"
		} else {
			feeUnitSource = "configured fee_unit"
		}
		c.options.logger.Info().Msgf("using fee unit: %s from %s", feeUnit, feeUnitSource)
	}
	return nil
}
