package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/logging"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/utils/must"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
)

type Cache interface {
	GetModel(ctx context.Context, key string, model interface{}) error
	SetModel(ctx context.Context, key string, model interface{}, ttl time.Duration, dependencies ...string) error
}

type cacheOptions struct {
	opts []cachestore.ClientOps
}

func (o *cacheOptions) configureLogger(cfg *config.AppConfig, logger zerolog.Logger) *cacheOptions {
	cachestoreLogger := logging.CreateGormLoggerAdapter(&logger, "cachestore")
	o.opts = append(o.opts, cachestore.WithLogger(cachestoreLogger))
	if logger.GetLevel() == zerolog.DebugLevel || logger.GetLevel() == zerolog.TraceLevel {
		o.opts = append(o.opts, cachestore.WithDebugging())
	}
	return o
}

func (o *cacheOptions) configureEngine(cfg *config.AppConfig) *cacheOptions {
	if cfg.Cache.Engine == cachestore.Redis {
		o.opts = append(o.opts, cachestore.WithRedis(&cachestore.RedisConfig{
			DependencyMode:        cfg.Cache.Redis.DependencyMode,
			MaxActiveConnections:  cfg.Cache.Redis.MaxActiveConnections,
			MaxConnectionLifetime: cfg.Cache.Redis.MaxConnectionLifetime,
			MaxIdleConnections:    cfg.Cache.Redis.MaxIdleConnections,
			MaxIdleTimeout:        cfg.Cache.Redis.MaxIdleTimeout,
			URL:                   cfg.Cache.Redis.URL,
			UseTLS:                cfg.Cache.Redis.UseTLS,
		}))
	} else if cfg.Cache.Engine == cachestore.FreeCache {
		o.opts = append(o.opts, cachestore.WithFreeCache())
	} else {
		panic(fmt.Sprintf("invalid configuration: unsupported cache engine: %s", cfg.Cache.Engine))
	}

	return o
}

func NewCache(cfg *config.AppConfig, logger zerolog.Logger) Cache {
	logger.With().Str("service", "cache").Logger()

	options := cacheOptions{make([]cachestore.ClientOps, 0)}

	options.configureLogger(cfg, logger).configureEngine(cfg)

	cache, err := cachestore.NewClient(context.Background(), options.opts...)
	must.HaveNoErrorf(err, "failed to create cache storage: %v", err)

	return cache
}
