package paymail_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	broadcast_client_mock "github.com/bitcoin-sv/go-broadcast-client/broadcast/broadcast-client-mock"
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	paymailclient "github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	xtester "github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/jarcoal/httpmock"
	"github.com/mrz1836/go-cache"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDomain    = "test.com"
	testServerURL = "https://" + testDomain + "/api/v1/" + paymail.DefaultServiceName
)

func Test_GetP2P(t *testing.T) {
	t.Parallel()

	t.Run("no p2p capabilities", func(t *testing.T) {
		client := paymailmock.MockClient(testDomain)
		client.WillRespondWithBasicCapabilities()

		paymailClient, err := paymailclient.NewServiceClient(nil, client)
		require.NoError(t, err)

		hasP2P, p2pDestinationURL, p2pSubmitTxURL, _ := paymailClient.GetP2P(context.Background(), testDomain)
		assert.False(t, hasP2P)
		assert.Equal(t, "", p2pDestinationURL)
		assert.Equal(t, "", p2pSubmitTxURL)
	})

	t.Run("valid p2p capabilities", func(t *testing.T) {
		client := paymailmock.MockClient(testDomain)
		client.WillRespondWithP2PCapabilities()

		paymailClient, err := paymailclient.NewServiceClient(nil, client)
		require.NoError(t, err)

		hasP2P, p2pDestinationURL, p2pSubmitTxURL, format := paymailClient.GetP2P(context.Background(), testDomain)
		assert.True(t, hasP2P)
		assert.Equal(t, paymailclient.BasicPaymailPayloadFormat, format)
		assert.Equal(t, client.GetMockedP2PPaymentDestinationURL(testDomain), p2pDestinationURL)
		assert.Equal(t, client.GetMockedP2PTransactionURL(testDomain), p2pSubmitTxURL)
	})

	t.Run("valid beef capabilities", func(t *testing.T) {
		client := paymailmock.MockClient(testDomain)
		client.WillRespondWithP2PWithBEEFCapabilities()

		paymailClient, err := paymailclient.NewServiceClient(nil, client)
		require.NoError(t, err)

		hasP2P, p2pDestinationURL, p2pSubmitTxURL, format := paymailClient.GetP2P(context.Background(), testDomain)
		assert.True(t, hasP2P)
		assert.Equal(t, paymailclient.BeefPaymailPayloadFormat, format)
		assert.Equal(t, client.GetMockedP2PPaymentDestinationURL(testDomain), p2pDestinationURL)
		assert.Equal(t, client.GetMockedBEEFTransactionURL(testDomain), p2pSubmitTxURL)
	})
}

func Test_StartP2PTransaction(t *testing.T) {
	// t.Parallel() mocking does not allow parallel tests

	const testAlias = "tester"

	t.Run("valid response", func(t *testing.T) {
		client := paymailmock.CreatePaymailClientService(testDomain)
		client.WillRespondWithP2PCapabilities()

		payload, err := client.StartP2PTransaction(
			testAlias,
			testDomain,
			client.GetMockedP2PPaymentDestinationURL(testDomain),
			1000,
		)
		require.NoError(t, err)
		require.NotNil(t, payload)
		assert.Equal(t, "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7", payload.Reference)
		assert.Equal(t, 1, len(payload.Outputs))
		assert.Equal(t, "16fkwYn8feXEbK7iCTg5KMx9Rx9GzZ9HuE", payload.Outputs[0].Address)
		assert.Equal(t, uint64(1000), payload.Outputs[0].Satoshis)
		assert.Equal(t, "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac", payload.Outputs[0].Script)
	})

	t.Run("error - address not found", func(t *testing.T) {
		client := paymailmock.CreatePaymailClientService(testDomain)
		client.WillRespondWithP2PCapabilities()
		client.WillRespondOnCapability(paymail.BRFCP2PPaymentDestination).WithNotFound()

		httpmock.RegisterResponder(http.MethodPost, testServerURL+"/p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusNotFound,
				`{"message": "not found"}`,
			),
		)

		payload, err := client.StartP2PTransaction(
			testAlias, testDomain,
			client.GetMockedP2PPaymentDestinationURL(testDomain), 1000,
		)
		require.Error(t, err)
		assert.Nil(t, payload)
	})
}

// Test_getCapabilities will test the method getCapabilities()
func Test_GetCapabilities(t *testing.T) {
	const (
		testIdleTimeout          = 240 * time.Second
		testMaxActiveConnections = 0
		testMaxConnLifetime      = 60 * time.Second
		testMaxIdleConnections   = 10
		testQueueName            = "test_queue"
		cacheKeyCapabilities     = "paymail-capabilities-"
	)

	// t.Parallel() mocking does not allow parallel tests
	bc := broadcast_client_mock.Builder().
		WithMockArc(broadcast_client_mock.MockSuccess).
		Build()

	t.Run("valid response - no cache found", func(t *testing.T) {
		client := paymailmock.MockClient(testDomain)
		client.WillRespondWithBasicCapabilities()

		redisClient, redisConn := xtester.LoadMockRedis(
			testIdleTimeout,
			testMaxConnLifetime,
			testMaxActiveConnections,
			testMaxIdleConnections,
		)
		logger := zerolog.Nop()

		tc, err := engine.NewClient(context.Background(),
			engine.WithRedisConnection(redisClient),
			engine.WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			engine.WithSQLite(&datastore.SQLiteConfig{Shared: true}),
			engine.WithChainstateOptions(false, false, false, false),
			engine.WithDebugging(),
			engine.WithBroadcastClient(bc),
			engine.WithLogger(&logger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer func() {
			closeClient(context.Background(), t, tc)
		}()

		// Get command
		getCmd := redisConn.Command(cache.GetCommand, cacheKeyCapabilities+testDomain).Expect(nil)

		paymailClient, err := paymailclient.NewServiceClient(tc.Cachestore(), client)
		require.NoError(t, err)

		var payload *paymail.CapabilitiesPayload
		payload, err = paymailClient.GetCapabilities(
			context.Background(), testDomain,
		)
		require.NoError(t, err)
		require.NotNil(t, payload)
		assert.Equal(t, true, getCmd.Called)
		assert.Equal(t, paymail.DefaultBsvAliasVersion, payload.BsvAlias)
		assert.Equal(t, 3, len(payload.Capabilities))
	})

	t.Run("[mocked] - server error", func(t *testing.T) {
		client := paymailmock.MockClient(testDomain)
		client.WillRespondWithErrorOnCapabilities()

		redisClient, redisConn := xtester.LoadMockRedis(
			testIdleTimeout,
			testMaxConnLifetime,
			testMaxActiveConnections,
			testMaxIdleConnections,
		)
		logger := zerolog.Nop()

		tc, err := engine.NewClient(context.Background(),
			engine.WithRedisConnection(redisClient),
			engine.WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			engine.WithSQLite(&datastore.SQLiteConfig{Shared: true}),
			engine.WithChainstateOptions(false, false, false, false),
			engine.WithDebugging(),
			engine.WithBroadcastClient(bc),
			engine.WithLogger(&logger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer func() {
			closeClient(context.Background(), t, tc)
		}()

		// Get command
		getCmd := redisConn.Command(cache.GetCommand, cacheKeyCapabilities+testDomain).Expect(nil)

		paymailClient, err := paymailclient.NewServiceClient(tc.Cachestore(), client)
		require.NoError(t, err)

		var payload *paymail.CapabilitiesPayload
		payload, err = paymailClient.GetCapabilities(
			context.Background(), testDomain,
		)
		require.Error(t, err)
		require.Nil(t, payload)
		assert.Equal(t, true, getCmd.Called)
	})

	t.Run("valid response - no cache found", func(t *testing.T) {
		client := paymailmock.CreatePaymailClientService(testDomain)
		client.WillRespondWithBasicCapabilities()

		logger := zerolog.Nop()
		tcOpts := defaultClientOpts(true, true)
		tcOpts = append(tcOpts, engine.WithLogger(&logger))

		tc, err := engine.NewClient(
			context.Background(),
			tcOpts...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer func() {
			closeClient(context.Background(), t, tc)
		}()

		var payload *paymail.CapabilitiesPayload
		payload, err = client.GetCapabilities(
			context.Background(), testDomain,
		)
		require.NoError(t, err)
		require.NotNil(t, payload)
		assert.Equal(t, paymail.DefaultBsvAliasVersion, payload.BsvAlias)
		assert.Equal(t, 3, len(payload.Capabilities))
	})

	t.Run("multiple requests for same capabilities", func(t *testing.T) {
		client := paymailmock.CreatePaymailClientService(testDomain)
		client.WillRespondWithBasicCapabilities()

		logger := zerolog.Nop()
		tcOpts := defaultClientOpts(true, true)
		tcOpts = append(tcOpts, engine.WithLogger(&logger))

		tc, err := engine.NewClient(
			context.Background(),
			tcOpts...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer func() {
			closeClient(context.Background(), t, tc)
		}()

		var payload *paymail.CapabilitiesPayload
		payload, err = client.GetCapabilities(
			context.Background(), testDomain,
		)
		require.NoError(t, err)
		require.NotNil(t, payload)
		assert.Equal(t, paymail.DefaultBsvAliasVersion, payload.BsvAlias)
		assert.Equal(t, 3, len(payload.Capabilities))

		time.Sleep(1 * time.Second)

		payload, err = client.GetCapabilities(
			context.Background(), testDomain,
		)
		require.NoError(t, err)
		require.NotNil(t, payload)
		assert.Equal(t, paymail.DefaultBsvAliasVersion, payload.BsvAlias)
		assert.Equal(t, 3, len(payload.Capabilities))
	})
}

func closeClient(ctx context.Context, t *testing.T, client engine.ClientInterface) {
	time.Sleep(1 * time.Second)
	require.NoError(t, client.Close(ctx))
}

// defaultClientOpts will return a default set of client options required to load the new client
func defaultClientOpts(debug, shared bool) []engine.ClientOps {
	tqc := taskmanager.DefaultTaskQConfig(xtester.RandomTablePrefix())
	tqc.MaxNumWorker = 2
	tqc.MaxNumFetcher = 2
	bc := broadcast_client_mock.Builder().
		WithMockArc(broadcast_client_mock.MockNilQueryTxResp).
		Build()

	opts := make([]engine.ClientOps, 0)
	opts = append(
		opts,
		engine.WithTaskqConfig(tqc),
		engine.WithSQLite(xtester.SQLiteTestConfig(debug, shared)),
		engine.WithChainstateOptions(false, false, false, false),
		engine.WithBroadcastClient(bc),
	)
	if debug {
		opts = append(opts, engine.WithDebugging())
	}

	return opts
}
