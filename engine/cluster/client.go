package cluster

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/logging"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type (

	// Client is the client (configuration)
	Client struct {
		pubSubService
		options *clientOptions
	}

	// clientOptions holds all the configuration for the client
	clientOptions struct {
		coordinator  Coordinator     // which coordinator to use, either 'memory' or 'redis'
		debug        bool            // For extra logs and additional debug information
		logger       *zerolog.Logger // Internal logger interface
		prefix       string          // the cluster key prefix to use before all keys
		redisOptions *redis.Options
	}
)

// NewClient create new cluster client
func NewClient(ctx context.Context, opts ...ClientOps) (*Client, error) {
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

	if client.options.coordinator == CoordinatorRedis {
		pubSubClient, err := NewRedisPubSub(ctx, client.options.redisOptions)
		if err != nil {
			return nil, err
		}
		pubSubClient.debug = client.IsDebug()
		pubSubClient.logger = client.options.logger
		pubSubClient.prefix = client.GetClusterPrefix()
		client.pubSubService = pubSubClient
	} else {
		pubSubClient, err := NewMemoryPubSub(ctx)
		if err != nil {
			return nil, err
		}

		pubSubClient.debug = client.IsDebug()
		pubSubClient.logger = client.options.logger
		pubSubClient.prefix = client.GetClusterPrefix()
		client.pubSubService = pubSubClient
	}

	// Return the client
	return client, nil
}

// IsDebug returns whether debugging is on or off
func (c *Client) IsDebug() bool {
	return c.options.debug
}

// GetClusterPrefix returns the cluster key prefix that can be used in things like Redis
func (c *Client) GetClusterPrefix() string {
	return c.options.prefix
}
