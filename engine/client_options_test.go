package engine

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/logging"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/coocood/freecache"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewRelicOptions will test the method enable()
func Test_newRelicOptions_enable(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("enable with valid app", func(t *testing.T) {
		app, err := tester.GetNewRelicApp(defaultNewRelicApp)
		require.NoError(t, err)
		require.NotNil(t, app)

		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithNewRelic(app))
		opts = append(opts, WithLogger(&testLogger))

		var tc ClientInterface
		tc, err = NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			opts...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		tc.EnableNewRelic()
		assert.Equal(t, true, tc.IsNewRelicEnabled())
	})

	t.Run("enable with invalid app", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithNewRelic(nil))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		tc.EnableNewRelic()
		assert.Equal(t, false, tc.IsNewRelicEnabled())
	})
}

// Test_newRelicOptions_getOrStartTxn will test the method getOrStartTxn()
func Test_newRelicOptions_getOrStartTxn(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("Get a valid ctx and txn", func(t *testing.T) {
		app, err := tester.GetNewRelicApp(defaultNewRelicApp)
		require.NoError(t, err)
		require.NotNil(t, app)

		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithNewRelic(app), WithLogger(&testLogger))

		var tc ClientInterface
		tc, err = NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			opts...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		ctx := tc.GetOrStartTxn(context.Background(), "test-name")
		assert.NotNil(t, ctx)

		txn := newrelic.FromContext(ctx)
		assert.NotNil(t, txn)
	})

	t.Run("invalid ctx and txn", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithNewRelic(nil), WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		ctx := tc.GetOrStartTxn(context.Background(), "test-name")
		assert.NotNil(t, ctx)

		txn := newrelic.FromContext(ctx)
		assert.Nil(t, txn)
	})
}

// TestClient_defaultModelOptions will test the method DefaultModelOptions()
func TestClient_defaultModelOptions(t *testing.T) {
	t.Parallel()

	t.Run("default options", func(t *testing.T) {
		dco := defaultClientOptions()
		require.NotNil(t, dco)

		require.NotNil(t, dco.cacheStore)
		require.Nil(t, dco.cacheStore.ClientInterface)
		require.NotNil(t, dco.cacheStore.options)
		assert.Equal(t, 0, len(dco.cacheStore.options))

		require.NotNil(t, dco.dataStore)
		require.Nil(t, dco.dataStore.ClientInterface)
		require.NotNil(t, dco.dataStore.options)
		assert.Equal(t, 1, len(dco.dataStore.options))

		require.NotNil(t, dco.newRelic)

		require.NotNil(t, dco.paymail)

		assert.Equal(t, defaultUserAgent, dco.userAgent)

		require.NotNil(t, dco.taskManager)

		assert.Nil(t, dco.logger)
	})
}

// TestWithUserAgent will test the method WithUserAgent()
func TestWithUserAgent(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithUserAgent("")
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("empty user agent", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithUserAgent(""), WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotEqual(t, "", tc.UserAgent())
		assert.Equal(t, defaultUserAgent, tc.UserAgent())
	})

	t.Run("custom user agent", func(t *testing.T) {
		customAgent := "custom-user-agent"

		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithUserAgent(customAgent), WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotEqual(t, defaultUserAgent, tc.UserAgent())
		assert.Equal(t, customAgent, tc.UserAgent())
	})
}

// TestWithNewRelic will test the method WithNewRelic()
func TestWithNewRelic(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithNewRelic(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})
}

// TestWithDebugging will test the method WithDebugging()
func TestWithDebugging(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithDebugging()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("set debug (with cache and Datastore)", func(t *testing.T) {
		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			DefaultClientOpts(true, true)...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, true, tc.IsDebug())
		assert.Equal(t, true, tc.Cachestore().IsDebug())
		assert.Equal(t, true, tc.Datastore().IsDebug())
	})
}

// TestWithEncryption will test the method WithEncryption()
func TestWithEncryption(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithEncryption("")
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("empty key", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithEncryption(""))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, false, tc.IsEncryptionKeySet())
	})

	t.Run("custom encryption key", func(t *testing.T) {
		key, _ := utils.RandomHex(32)
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithEncryption(key))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, true, tc.IsEncryptionKeySet())
	})
}

// TestWithRedis will test the method WithRedis()
func TestWithRedis(t *testing.T) {
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithRedis(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("using valid config", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping live local redis tests")
		}

		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix())),
			WithRedis(&cachestore.RedisConfig{
				URL: cachestore.RedisPrefix + "localhost:6379",
			}),
			WithSQLite(tester.SQLiteTestConfig(false, true)),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&testLogger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		cs := tc.Cachestore()
		require.NotNil(t, cs)
		assert.Equal(t, cachestore.Redis, cs.Engine())
	})

	t.Run("missing redis prefix", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping live local redis tests")
		}

		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix())),
			WithRedis(&cachestore.RedisConfig{
				URL: "localhost:6379",
			}),
			WithSQLite(tester.SQLiteTestConfig(false, true)),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&testLogger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		cs := tc.Cachestore()
		require.NotNil(t, cs)
		assert.Equal(t, cachestore.Redis, cs.Engine())
	})
}

// TestWithRedisConnection will test the method WithRedisConnection()
func TestWithRedisConnection(t *testing.T) {
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithRedisConnection(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("using a nil connection", func(t *testing.T) {
		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix())),
			WithRedisConnection(nil),
			WithSQLite(tester.SQLiteTestConfig(false, true)),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&testLogger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		cs := tc.Cachestore()
		require.NotNil(t, cs)
		assert.Equal(t, cachestore.FreeCache, cs.Engine())
	})

	t.Run("using an existing connection", func(t *testing.T) {
		client, conn := tester.LoadMockRedis(10*time.Second, 10*time.Second, 10, 10)
		require.NotNil(t, client)
		require.NotNil(t, conn)

		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix())),
			WithRedisConnection(client),
			WithSQLite(tester.SQLiteTestConfig(false, true)),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&testLogger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		cs := tc.Cachestore()
		require.NotNil(t, cs)
		assert.Equal(t, cachestore.Redis, cs.Engine())
	})
}

// TestWithFreeCache will test the method WithFreeCache()
func TestWithFreeCache(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithFreeCache()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("using FreeCache", func(t *testing.T) {
		testLogger := zerolog.Nop()
		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			WithFreeCache(),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			WithSQLite(&datastore.SQLiteConfig{Shared: true}),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&testLogger))
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		cs := tc.Cachestore()
		require.NotNil(t, cs)
		assert.Equal(t, cachestore.FreeCache, cs.Engine())
	})
}

// TestWithFreeCacheConnection will test the method WithFreeCacheConnection()
func TestWithFreeCacheConnection(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithFreeCacheConnection(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("using a nil client", func(t *testing.T) {
		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			WithFreeCacheConnection(nil),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			WithSQLite(&datastore.SQLiteConfig{Shared: true}),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&testLogger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)

		defer CloseClient(context.Background(), t, tc)

		cs := tc.Cachestore()
		require.NotNil(t, cs)
		assert.Equal(t, cachestore.FreeCache, cs.Engine())
	})

	t.Run("using an existing connection", func(t *testing.T) {
		fc := freecache.NewCache(cachestore.DefaultCacheSize)
		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			WithFreeCacheConnection(fc),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			WithSQLite(&datastore.SQLiteConfig{Shared: true}),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&testLogger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		cs := tc.Cachestore()
		require.NotNil(t, cs)
		assert.Equal(t, cachestore.FreeCache, cs.Engine())
	})
}

// TestWithPaymailClient will test the method WithPaymailClient()
func TestWithPaymailClient(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("using a nil driver, automatically makes paymail client", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithPaymailClient(nil))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.PaymailClient())
	})

	t.Run("custom paymail client", func(t *testing.T) {
		p, err := paymail.NewClient()
		require.NoError(t, err)
		require.NotNil(t, p)

		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithPaymailClient(p))
		opts = append(opts, WithLogger(&testLogger))

		var tc ClientInterface
		tc, err = NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.PaymailClient())
		assert.Equal(t, p, tc.PaymailClient())
	})
}

// TestWithTaskQ will test the method WithTaskQ()
func TestWithTaskQ(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	// todo: test cases where config is nil, or cannot load TaskQ

	t.Run("using taskq using memory", func(t *testing.T) {
		tcOpts := DefaultClientOpts(true, true)
		tcOpts = append(tcOpts, WithLogger(&testLogger))

		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			tcOpts...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		tm := tc.Taskmanager()
		require.NotNil(t, tm)
		assert.Equal(t, taskmanager.FactoryMemory, tm.Factory())
	})

	t.Run("using taskq using redis", func(t *testing.T) {
		if testing.Short() {
			t.Skip("skipping live local redis tests")
		}

		tc, err := NewClient(
			tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx),
			WithTaskqConfig(
				taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix(), taskmanager.WithRedis("localhost:6379")),
			),
			WithRedis(&cachestore.RedisConfig{
				URL: cachestore.RedisPrefix + "localhost:6379",
			}),
			WithSQLite(tester.SQLiteTestConfig(false, true)),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&testLogger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		tm := tc.Taskmanager()
		require.NotNil(t, tm)
		assert.Equal(t, taskmanager.FactoryRedis, tm.Factory())
	})
}

// TestWithLogger will test the method WithLogger()
func TestWithLogger(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithLogger(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithLogger(nil))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.Logger())
		assert.Equal(t, logging.GetDefaultLogger(), tc.Logger())
	})

	t.Run("test applying option", func(t *testing.T) {
		customLogger := zerolog.Nop()
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithLogger(&customLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, &customLogger, tc.Logger())
	})
}

// TestWithModels will test the method WithModels()
func TestWithModels(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithModels()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("empty models - returns default models", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithModels())
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, []string{
			ModelXPub.String(), ModelAccessKey.String(),
			ModelDraftTransaction.String(), ModelTransaction.String(),
			ModelSyncTransaction.String(), ModelDestination.String(),
			ModelUtxo.String(),
		}, tc.GetModelNames())
	})

	t.Run("add custom models", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithModels(newPaymail(testPaymail)))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, []string{
			ModelXPub.String(), ModelAccessKey.String(),
			ModelDraftTransaction.String(), ModelTransaction.String(),
			ModelSyncTransaction.String(), ModelDestination.String(),
			ModelUtxo.String(), ModelPaymailAddress.String(),
		}, tc.GetModelNames())
	})
}

// TestWithIUCDisabled will test the method WithIUCDisabled()
func TestWithIUCDisabled(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithIUCDisabled()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("default options", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, true, tc.IsIUCEnabled())
	})

	t.Run("iuc disabled", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithIUCDisabled())
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, false, tc.IsIUCEnabled())
	})
}

// TestWithHTTPClient will test the method WithHTTPClient()
func TestWithHTTPClient(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithHTTPClient(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithHTTPClient(nil))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.HTTPClient())
	})

	t.Run("test applying option", func(t *testing.T) {
		customClient := &http.Client{}
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithHTTPClient(customClient))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, customClient, tc.HTTPClient())
	})
}

// TestWithCustomCachestore will test the method WithCustomCachestore()
func TestWithCustomCachestore(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithCustomCachestore(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithCustomCachestore(nil))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.Cachestore())
	})

	t.Run("test applying option", func(t *testing.T) {
		customCache, err := cachestore.NewClient(context.Background())
		require.NoError(t, err)

		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithCustomCachestore(customCache))
		opts = append(opts, WithLogger(&testLogger))

		var tc ClientInterface
		tc, err = NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, customCache, tc.Cachestore())
	})
}

// TestWithCustomDatastore will test the method WithCustomDatastore()
func TestWithCustomDatastore(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithCustomDatastore(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithCustomDatastore(nil))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.Datastore())
	})

	t.Run("test applying option", func(t *testing.T) {
		customData, err := datastore.NewClient(context.Background())
		require.NoError(t, err)

		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithCustomDatastore(customData))
		opts = append(opts, WithLogger(&testLogger))

		var tc ClientInterface
		tc, err = NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, customData, tc.Datastore())
	})

	// Attempt to remove a file created during the test
	t.Cleanup(func() {
		_ = os.Remove("datastore.db")
	})
}

// TestWithAutoMigrate will test the method WithAutoMigrate()
func TestWithAutoMigrate(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithAutoMigrate()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("no additional models, just base models", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithAutoMigrate())
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, []string{
			ModelXPub.String(),
			ModelAccessKey.String(),
			ModelDraftTransaction.String(),
			ModelTransaction.String(),
			ModelSyncTransaction.String(),
			ModelDestination.String(),
			ModelUtxo.String(),
		}, tc.GetModelNames())
	})

	t.Run("one additional model", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithAutoMigrate(newPaymail(testPaymail)))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, []string{
			ModelXPub.String(),
			ModelAccessKey.String(),
			ModelDraftTransaction.String(),
			ModelTransaction.String(),
			ModelSyncTransaction.String(),
			ModelDestination.String(),
			ModelUtxo.String(),
			ModelPaymailAddress.String(),
		}, tc.GetModelNames())
	})
}

// TestWithMigrationDisabled will test the method WithMigrationDisabled()
func TestWithMigrationDisabled(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithMigrationDisabled()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("default options", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, true, tc.IsMigrationEnabled())
	})

	t.Run("migration disabled", func(t *testing.T) {
		opts := DefaultClientOpts(false, true)
		opts = append(opts, WithMigrationDisabled())
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(tester.GetNewRelicCtx(t, defaultNewRelicApp, defaultNewRelicTx), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, false, tc.IsMigrationEnabled())
	})
}
