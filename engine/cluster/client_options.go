package cluster

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const clientOptsPrefix = "bsv_"

// ClientOps allow functional options to be supplied
// that overwrite default client options.
type ClientOps func(c *clientOptions)

// defaultClientOptions will return an clientOptions struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultClientOptions() *clientOptions {
	// Set the default options
	return &clientOptions{
		debug:       false,
		coordinator: CoordinatorMemory,
		prefix:      clientOptsPrefix,
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

// WithRedis will enable redis cluster coordinator
func WithRedis(redisOptions *redis.Options) ClientOps {
	return func(c *clientOptions) {
		c.coordinator = CoordinatorRedis
		c.redisOptions = redisOptions
	}
}

// WithKeyPrefix will set the prefix to use for all keys in the cluster coordinator
func WithKeyPrefix(prefix string) ClientOps {
	return func(c *clientOptions) {
		c.prefix = prefix
	}
}
