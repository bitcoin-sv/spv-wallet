package engine

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/logging"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/coocood/freecache"
	"github.com/mrz1836/go-cachestore"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

		require.NotNil(t, dco.paymail)

		assert.Equal(t, defaultUserAgent, dco.userAgent)

		require.NotNil(t, dco.taskManager)

		assert.Nil(t, dco.logger)
	})
}

func TestWithUserAgent(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithUserAgent("")
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("empty user agent", func(t *testing.T) {
		opts := DefaultClientOpts()
		opts = append(opts, WithUserAgent(""), WithLogger(&testLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotEqual(t, "", tc.UserAgent())
		assert.Equal(t, defaultUserAgent, tc.UserAgent())
	})

	t.Run("custom user agent", func(t *testing.T) {
		customAgent := "custom-user-agent"

		opts := DefaultClientOpts()
		opts = append(opts, WithUserAgent(customAgent), WithLogger(&testLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotEqual(t, defaultUserAgent, tc.UserAgent())
		assert.Equal(t, customAgent, tc.UserAgent())
	})
}

func TestWithDebugging(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithDebugging()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("set debug (with cache and Datastore)", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			append(DefaultClientOpts(), WithDebugging())...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, true, tc.Cachestore().IsDebug())
		assert.Equal(t, true, tc.Datastore().IsDebug())
	})
}

func TestWithEncryption(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithEncryption("")
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("empty key", func(t *testing.T) {
		opts := DefaultClientOpts()
		opts = append(opts, WithEncryption(""))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, false, tc.IsEncryptionKeySet())
	})

	t.Run("custom encryption key", func(t *testing.T) {
		key, _ := utils.RandomHex(32)
		opts := DefaultClientOpts()
		opts = append(opts, WithEncryption(key))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, true, tc.IsEncryptionKeySet())
	})
}

func TestWithRedisConnection(t *testing.T) {
	testLogger := zerolog.Nop()

	t.Run("using a nil connection", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix())),
			WithRedisConnection(nil),
			WithSQLite(tester.SQLiteTestConfig()),
			WithCustomFeeUnit(mockFeeUnit),
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
			context.Background(),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(tester.RandomTablePrefix())),
			WithRedisConnection(client),
			WithSQLite(tester.SQLiteTestConfig()),
			WithCustomFeeUnit(mockFeeUnit),
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

func TestWithFreeCache(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithFreeCache()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("using FreeCache", func(t *testing.T) {
		testLogger := zerolog.Nop()
		tc, err := NewClient(
			context.Background(),
			WithFreeCache(),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			WithSQLite(tester.SQLiteTestConfig()),
			WithCustomFeeUnit(mockFeeUnit),
			WithLogger(&testLogger))
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		cs := tc.Cachestore()
		require.NotNil(t, cs)
		assert.Equal(t, cachestore.FreeCache, cs.Engine())
	})
}

func TestWithFreeCacheConnection(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithFreeCacheConnection(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("using a nil client", func(t *testing.T) {
		tc, err := NewClient(
			context.Background(),
			WithFreeCacheConnection(nil),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			WithSQLite(tester.SQLiteTestConfig()),
			WithCustomFeeUnit(mockFeeUnit),
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
			context.Background(),
			WithFreeCacheConnection(fc),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			WithSQLite(tester.SQLiteTestConfig()),
			WithCustomFeeUnit(mockFeeUnit),
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

func TestWithPaymailClient(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("using a nil driver, automatically makes paymail client", func(t *testing.T) {
		opts := DefaultClientOpts()
		opts = append(opts, WithPaymailClient(nil))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.PaymailClient())
	})

	t.Run("custom paymail client", func(t *testing.T) {
		p, err := paymail.NewClient()
		require.NoError(t, err)
		require.NotNil(t, p)

		opts := DefaultClientOpts()
		opts = append(opts, WithPaymailClient(p))
		opts = append(opts, WithLogger(&testLogger))

		var tc ClientInterface
		tc, err = NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.PaymailClient())
		assert.Equal(t, p, tc.PaymailClient())
	})
}

func TestWithTaskQ(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	// todo: test cases where config is nil, or cannot load TaskQ

	t.Run("using taskq using memory", func(t *testing.T) {
		tcOpts := DefaultClientOpts()
		tcOpts = append(tcOpts, WithLogger(&testLogger))

		tc, err := NewClient(
			context.Background(),
			tcOpts...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		tm := tc.Taskmanager()
		require.NotNil(t, tm)
		assert.Equal(t, taskmanager.FactoryMemory, tm.Factory())
	})
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	t.Run("check type", func(t *testing.T) {
		opt := WithLogger(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		opts := DefaultClientOpts()
		opts = append(opts, WithLogger(nil))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.Logger())
		assert.Equal(t, logging.GetDefaultLogger(), tc.Logger())
	})

	t.Run("test applying option", func(t *testing.T) {
		customLogger := zerolog.Nop()
		opts := DefaultClientOpts()
		opts = append(opts, WithLogger(&customLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, &customLogger, tc.Logger())
	})
}

func TestWithIUCDisabled(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithIUCDisabled()
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("default options", func(t *testing.T) {
		opts := DefaultClientOpts()
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, true, tc.IsIUCEnabled())
	})

	t.Run("iuc disabled", func(t *testing.T) {
		opts := DefaultClientOpts()
		opts = append(opts, WithIUCDisabled())
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, false, tc.IsIUCEnabled())
	})
}

func TestWithCustomCachestore(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithCustomCachestore(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		opts := DefaultClientOpts()
		opts = append(opts, WithCustomCachestore(nil))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.Cachestore())
	})

	t.Run("test applying option", func(t *testing.T) {
		customCache, err := cachestore.NewClient(context.Background())
		require.NoError(t, err)

		opts := DefaultClientOpts()
		opts = append(opts, WithCustomCachestore(customCache))
		opts = append(opts, WithLogger(&testLogger))

		var tc ClientInterface
		tc, err = NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.Equal(t, customCache, tc.Cachestore())
	})
}

func TestWithCustomDatastore(t *testing.T) {
	t.Parallel()
	testLogger := zerolog.Nop()

	t.Run("check type", func(t *testing.T) {
		opt := WithCustomDatastore(nil)
		assert.IsType(t, *new(ClientOps), opt)
	})

	t.Run("test applying nil", func(t *testing.T) {
		opts := DefaultClientOpts()
		opts = append(opts, WithCustomDatastore(nil))
		opts = append(opts, WithLogger(&testLogger))

		tc, err := NewClient(context.Background(), opts...)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer CloseClient(context.Background(), t, tc)

		assert.NotNil(t, tc.Datastore())
	})

	t.Run("test applying option", func(t *testing.T) {
		customData, err := datastore.NewClient()
		require.NoError(t, err)

		opts := DefaultClientOpts()
		opts = append(opts, WithCustomDatastore(customData))
		opts = append(opts, WithLogger(&testLogger))

		var tc ClientInterface
		tc, err = NewClient(context.Background(), opts...)
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
