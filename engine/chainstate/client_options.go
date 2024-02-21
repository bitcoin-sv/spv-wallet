package chainstate

import (
	"context"
	"time"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/tonicpow/go-minercraft/v2"
)

// ClientOps allow functional options to be supplied
// that overwrite default client options.
type ClientOps func(c *clientOptions)

// defaultClientOptions will return an clientOptions struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultClientOptions() *clientOptions {
	// Set the default options
	return &clientOptions{
		config: &syncConfig{
			httpClient:       nil,
			minercraftConfig: defaultMinecraftConfig(),
			minercraft:       nil,
			network:          MainNet,
			queryTimeout:     defaultQueryTimeOut,
			broadcastClient:  nil,
			feeQuotes:        true,
			feeUnit:          nil, // fee has to be set explicitly or via fee quotes
		},
		debug:           false,
		newRelicEnabled: false,
		metrics:         nil,
	}
}

// getTxnCtx will check for an existing transaction
func (c *clientOptions) getTxnCtx(ctx context.Context) context.Context {
	if c.newRelicEnabled {
		txn := newrelic.FromContext(ctx)
		if txn != nil {
			ctx = newrelic.NewContext(ctx, txn)
		}
	}
	return ctx
}

// WithNewRelic will enable the NewRelic wrapper
func WithNewRelic() ClientOps {
	return func(c *clientOptions) {
		c.newRelicEnabled = true
	}
}

// WithDebugging will enable debugging mode
func WithDebugging() ClientOps {
	return func(c *clientOptions) {
		c.debug = true
	}
}

// WithHTTPClient will set a custom HTTP client
func WithHTTPClient(client HTTPInterface) ClientOps {
	return func(c *clientOptions) {
		if client != nil {
			c.config.httpClient = client
		}
	}
}

// WithMinercraft will set a custom Minercraft client
func WithMinercraft(client minercraft.ClientInterface) ClientOps {
	return func(c *clientOptions) {
		if client != nil {
			c.config.minercraft = client
		}
	}
}

// WithMAPI will specify mAPI as an API for minercraft client
func WithMAPI() ClientOps {
	return func(c *clientOptions) {
		c.config.minercraftConfig.apiType = minercraft.MAPI
	}
}

// WithQueryTimeout will set a different timeout for transaction querying
func WithQueryTimeout(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		if timeout > 0 {
			c.config.queryTimeout = timeout
		}
	}
}

// WithUserAgent will set the custom user agent
func WithUserAgent(agent string) ClientOps {
	return func(c *clientOptions) {
		if len(agent) > 0 {
			c.userAgent = agent
		}
	}
}

// WithNetwork will set the network to use
func WithNetwork(network Network) ClientOps {
	return func(c *clientOptions) {
		if len(network) > 0 {
			c.config.network = network
		}
	}
}

// WithLogger will set a custom logger
func WithLogger(customLogger *zerolog.Logger) ClientOps {
	return func(c *clientOptions) {
		if customLogger != nil {
			c.logger = customLogger
		}
	}
}

// WithExcludedProviders will set a list of excluded providers
func WithExcludedProviders(providers []string) ClientOps {
	return func(c *clientOptions) {
		if len(providers) > 0 {
			c.config.excludedProviders = providers
		}
	}
}

// WithFeeQuotes will set minercraftFeeQuotes flag as true
func WithFeeQuotes(enabled bool) ClientOps {
	return func(c *clientOptions) {
		c.config.feeQuotes = enabled
	}
}

// WithFeeUnit will set the fee unit
func WithFeeUnit(feeUnit *utils.FeeUnit) ClientOps {
	return func(c *clientOptions) {
		c.config.feeUnit = feeUnit
	}
}

// WithMinercraftAPIs will set miners APIs
func WithMinercraftAPIs(apis []*minercraft.MinerAPIs) ClientOps {
	return func(c *clientOptions) {
		c.config.minercraftConfig.minerAPIs = apis
	}
}

// WithBroadcastClient will set broadcast client APIs
func WithBroadcastClient(client broadcast.Client) ClientOps {
	return func(c *clientOptions) {
		c.config.broadcastClient = client
	}
}

// WithConnectionToBlockHeaderService will set Block Headers Service API settings.
func WithConnectionToBlockHeaderService(url, authToken string) ClientOps {
	return func(c *clientOptions) {
		c.config.blockHedersServiceClient = newBlockHeaderServiceClientProvider(url, authToken)
	}
}

// WithCallback will set broadcast callback settings
func WithCallback(callbackURL, callbackAuthToken string) ClientOps {
	return func(c *clientOptions) {
		c.config.callbackURL = callbackURL
		c.config.callbackToken = callbackAuthToken
	}
}

// WithMetrics will set metrics
func WithMetrics(metrics *metrics.Metrics) ClientOps {
	return func(c *clientOptions) {
		c.metrics = metrics
	}
}
