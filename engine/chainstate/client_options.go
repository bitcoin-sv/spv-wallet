package chainstate

import (
	"time"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/rs/zerolog"
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
			httpClient:            nil,
			broadcastClientConfig: defaultArcConfig(),
			network:               MainNet,
			queryTimeout:          defaultQueryTimeOut,
			broadcastClient:       nil,
			feeQuotes:             true,
			feeUnit:               nil, // fee has to be set explicitly or via fee quotes
		},
		debug:   false,
		metrics: nil,
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

// WithFeeQuotes will set feeQuotes flag as true
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

// WithBroadcastClient will set broadcast client APIs
func WithBroadcastClient(client broadcast.Client) ClientOps {
	return func(c *clientOptions) {
		c.config.broadcastClient = client
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
