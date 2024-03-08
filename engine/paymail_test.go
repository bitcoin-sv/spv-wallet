package engine

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/spv-wallet/engine/chainstate"
	"github.com/bitcoin-sv/spv-wallet/engine/taskmanager"
	xtester "github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/jarcoal/httpmock"
	"github.com/mrz1836/go-cache"
	"github.com/mrz1836/go-datastore"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAlias     = "tester"
	testDomain    = "test.com"
	testServerURL = "https://" + testDomain + "/api/v1/" + paymail.DefaultServiceName
	testOutput    = "76a9147f11c8f67a2781df0400ebfb1f31b4c72a780b9d88ac"
)

// newTestPaymailClient will return a client for testing purposes
func newTestPaymailClient(t *testing.T, domains []string) paymail.ClientInterface {
	newClient, err := xtester.PaymailMockClient(domains)
	require.NotNil(t, newClient)
	require.NoError(t, err)
	return newClient
}

// newTestPaymailConfig loads a basic test configuration
func newTestPaymailConfig(t *testing.T, domain string) *server.Configuration {
	pl := &server.PaymailServiceLocator{}
	pl.RegisterPaymailService(new(mockServiceProvider))

	c, err := server.NewConfig(
		pl,
		server.WithDomain(domain),
		server.WithP2PCapabilities(),
	)
	require.NoError(t, err)
	require.NotNil(t, c)
	return c
}

// mockValidResponse is used for mocking the response
func mockValidResponse(statusCode int, p2p bool, domain string) {
	httpmock.Reset()
	serverURL := "https://" + domain + "/api/v1/" + paymail.DefaultServiceName

	// Basic address resolution vs P2P
	if !p2p {
		httpmock.RegisterResponder(http.MethodGet, "https://"+domain+":443/.well-known/"+paymail.DefaultServiceName,
			httpmock.NewStringResponder(
				statusCode,
				`{"`+paymail.DefaultServiceName+`": "`+paymail.DefaultBsvAliasVersion+`","capabilities":{
"`+paymail.BRFCSenderValidation+`": false,
"`+paymail.BRFCPki+`": "`+serverURL+`/id/{alias}@{domain.tld}",
"`+paymail.BRFCPaymentDestination+`": "`+serverURL+`/address/{alias}@{domain.tld}"}
}`,
			),
		)
	} else {
		httpmock.RegisterResponder(http.MethodGet, "https://"+domain+":443/.well-known/"+paymail.DefaultServiceName,
			httpmock.NewStringResponder(
				statusCode,
				`{"`+paymail.DefaultServiceName+`": "`+paymail.DefaultBsvAliasVersion+`","capabilities":{
"`+paymail.BRFCSenderValidation+`": false,
"`+paymail.BRFCPki+`": "`+serverURL+`/id/{alias}@{domain.tld}",
"`+paymail.BRFCPaymentDestination+`": "`+serverURL+`/address/{alias}@{domain.tld}",
"`+paymail.BRFCP2PTransactions+`": "`+serverURL+`/receive-transaction/{alias}@{domain.tld}",
"`+paymail.BRFCP2PPaymentDestination+`": "`+serverURL+`/p2p-payment-destination/{alias}@{domain.tld}"}
}`,
			),
		)
	}

	httpmock.RegisterResponder(http.MethodPost, serverURL+"/p2p-payment-destination/"+testAlias+"@"+testDomain,
		httpmock.NewStringResponder(
			statusCode,
			`{"outputs": [{"script": "76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac","satoshis": 1000}],"reference": "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"}`,
		),
	)

	httpmock.RegisterResponder(http.MethodPost, serverURL+"/address/"+testAlias+"@"+domain,
		httpmock.NewStringResponder(
			statusCode,
			`{"output": "`+testOutput+`"}`,
		),
	)
}

// TestPaymailClient will test various Paymail client methods
func TestPaymailClient(t *testing.T) {
	t.Parallel()

	config := newTestPaymailConfig(t, testDomain)
	require.NotNil(t, config)

	client := newTestPaymailClient(t, []string{testDomain})
	require.NotNil(t, client)
}

func mockCapabilities(t *testing.T, p2pEnabled, beefEnabled bool) *paymail.CapabilitiesPayload {
	options := []server.ConfigOps{
		server.WithDomain("test.com"),
	}
	if p2pEnabled {
		options = append(options, server.WithP2PCapabilities())
	}
	if beefEnabled {
		options = append(options, server.WithBeefCapabilities())
	}

	pl := &server.PaymailServiceLocator{}
	pl.RegisterPaymailService(new(mockServiceProvider))

	config, err := server.NewConfig(pl, options...)
	assert.NoError(t, err)
	capPayload, err := config.EnrichCapabilities("test.com")
	if err != nil {
		return &paymail.CapabilitiesPayload{
			BsvAlias: "",
		}
	}
	return capPayload
}

// Test_hasP2P will test the method hasP2P()
func Test_hasP2P(t *testing.T) {
	t.Parallel()

	t.Run("no p2p capabilities", func(t *testing.T) {
		capabilities := mockCapabilities(t, false, false)
		success, p2pDestinationURL, p2pSubmitTxURL, _ := hasP2P(capabilities)
		assert.Equal(t, false, success)
		assert.Equal(t, "", p2pDestinationURL)
		assert.Equal(t, "", p2pSubmitTxURL)
	})

	t.Run("valid p2p capabilities", func(t *testing.T) {
		capabilities := mockCapabilities(t, true, false)

		success, p2pDestinationURL, p2pSubmitTxURL, _ := hasP2P(capabilities)
		assert.Equal(t, true, success)
		assert.Equal(t, capabilities.Capabilities[paymail.BRFCP2PPaymentDestination], p2pDestinationURL)
		assert.Equal(t, capabilities.Capabilities[paymail.BRFCP2PTransactions], p2pSubmitTxURL)
	})
}

// Test_hasP2P_beefCapabilities will test the method hasP2P() but with BEEF capabilities
func Test_hasP2P_beefCapabilities(t *testing.T) {
	t.Parallel()

	t.Run("no beef capabilities", func(t *testing.T) {
		capabilities := mockCapabilities(t, false, false)
		success, p2pDestinationURL, p2pSubmitTxURL, format := hasP2P(capabilities)
		assert.Equal(t, false, success)
		assert.Equal(t, BasicPaymailPayloadFormat, format)
		assert.Equal(t, "", p2pDestinationURL)
		assert.Equal(t, "", p2pSubmitTxURL)
	})

	t.Run("valid beef capabilities", func(t *testing.T) {
		capabilities := mockCapabilities(t, true, true)
		success, p2pDestinationURL, p2pSubmitTxURL, format := hasP2P(capabilities)
		assert.Equal(t, true, success)
		assert.Equal(t, BeefPaymailPayloadFormat, format)
		assert.Equal(t, capabilities.Capabilities[paymail.BRFCP2PPaymentDestination], p2pDestinationURL)
		assert.Equal(t, capabilities.Capabilities[paymail.BRFCBeefTransaction], p2pSubmitTxURL)
	})
}

// Test_startP2PTransaction will test the method startP2PTransaction()
func Test_startP2PTransaction(t *testing.T) {
	// t.Parallel() mocking does not allow parallel tests

	t.Run("[mocked] - valid response", func(t *testing.T) {
		client := newTestPaymailClient(t, []string{testDomain})

		mockValidResponse(http.StatusOK, true, testDomain)

		payload, err := startP2PTransaction(
			client, testAlias, testDomain,
			testServerURL+"/p2p-payment-destination/{alias}@{domain.tld}", 1000,
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
		client := newTestPaymailClient(t, []string{testDomain})

		httpmock.Reset()
		httpmock.RegisterResponder(http.MethodPost, testServerURL+"/p2p-payment-destination/"+testAlias+"@"+testDomain,
			httpmock.NewStringResponder(
				http.StatusNotFound,
				`{"message": "not found"}`,
			),
		)

		payload, err := startP2PTransaction(
			client, testAlias, testDomain,
			testServerURL+"/p2p-payment-destination/{alias}@{domain.tld}", 1000,
		)

		require.Error(t, err)
		assert.Nil(t, payload)
	})
}

// Test_getCapabilities will test the method getCapabilities()
func Test_getCapabilities(t *testing.T) {
	// t.Parallel() mocking does not allow parallel tests

	t.Run("[mocked] - valid response - no cache found", func(t *testing.T) {
		client := newTestPaymailClient(t, []string{testDomain})

		redisClient, redisConn := xtester.LoadMockRedis(
			testIdleTimeout,
			testMaxConnLifetime,
			testMaxActiveConnections,
			testMaxIdleConnections,
		)
		logger := zerolog.Nop()

		tc, err := NewClient(context.Background(),
			WithRedisConnection(redisClient),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			WithSQLite(&datastore.SQLiteConfig{Shared: true}),
			WithChainstateOptions(false, false, false, false),
			WithDebugging(),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&logger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer func() {
			time.Sleep(1 * time.Second)
			CloseClient(context.Background(), t, tc)
		}()

		// Get command
		getCmd := redisConn.Command(cache.GetCommand, cacheKeyCapabilities+testDomain).Expect(nil)

		mockValidResponse(http.StatusOK, false, testDomain)
		var payload *paymail.CapabilitiesPayload
		payload, err = getCapabilities(
			context.Background(), tc.Cachestore(), client, testDomain,
		)
		require.NoError(t, err)
		require.NotNil(t, payload)
		assert.Equal(t, true, getCmd.Called)
		assert.Equal(t, paymail.DefaultBsvAliasVersion, payload.BsvAlias)
		assert.Equal(t, 3, len(payload.Capabilities))
	})

	t.Run("[mocked] - server error", func(t *testing.T) {
		client := newTestPaymailClient(t, []string{testDomain})

		redisClient, redisConn := xtester.LoadMockRedis(
			testIdleTimeout,
			testMaxConnLifetime,
			testMaxActiveConnections,
			testMaxIdleConnections,
		)
		logger := zerolog.Nop()

		tc, err := NewClient(context.Background(),
			WithRedisConnection(redisClient),
			WithTaskqConfig(taskmanager.DefaultTaskQConfig(testQueueName)),
			WithSQLite(&datastore.SQLiteConfig{Shared: true}),
			WithChainstateOptions(false, false, false, false),
			WithDebugging(),
			WithMinercraft(&chainstate.MinerCraftBase{}),
			WithLogger(&logger),
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer func() {
			time.Sleep(1 * time.Second)
			CloseClient(context.Background(), t, tc)
		}()

		// Get command
		getCmd := redisConn.Command(cache.GetCommand, cacheKeyCapabilities+testDomain).Expect(nil)

		mockValidResponse(http.StatusBadRequest, false, testDomain)
		var payload *paymail.CapabilitiesPayload
		payload, err = getCapabilities(
			context.Background(), tc.Cachestore(), client, testDomain,
		)
		require.Error(t, err)
		require.Nil(t, payload)
		assert.Equal(t, true, getCmd.Called)
	})

	t.Run("valid response - no cache found", func(t *testing.T) {
		client := newTestPaymailClient(t, []string{testDomain})

		logger := zerolog.Nop()
		tcOpts := DefaultClientOpts(true, true)
		tcOpts = append(tcOpts, WithLogger(&logger))

		tc, err := NewClient(
			context.Background(),
			tcOpts...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer func() {
			time.Sleep(1 * time.Second)
			CloseClient(context.Background(), t, tc)
		}()

		mockValidResponse(http.StatusOK, false, testDomain)
		var payload *paymail.CapabilitiesPayload
		payload, err = getCapabilities(
			context.Background(), tc.Cachestore(), client, testDomain,
		)
		require.NoError(t, err)
		require.NotNil(t, payload)
		assert.Equal(t, paymail.DefaultBsvAliasVersion, payload.BsvAlias)
		assert.Equal(t, 3, len(payload.Capabilities))
	})

	t.Run("multiple requests for same capabilities", func(t *testing.T) {
		client := newTestPaymailClient(t, []string{testDomain})

		logger := zerolog.Nop()
		tcOpts := DefaultClientOpts(true, true)
		tcOpts = append(tcOpts, WithLogger(&logger))

		tc, err := NewClient(
			context.Background(),
			tcOpts...,
		)
		require.NoError(t, err)
		require.NotNil(t, tc)
		defer func() {
			time.Sleep(1 * time.Second)
			CloseClient(context.Background(), t, tc)
		}()

		mockValidResponse(http.StatusOK, false, testDomain)
		var payload *paymail.CapabilitiesPayload
		payload, err = getCapabilities(
			context.Background(), tc.Cachestore(), client, testDomain,
		)
		require.NoError(t, err)
		require.NotNil(t, payload)
		assert.Equal(t, paymail.DefaultBsvAliasVersion, payload.BsvAlias)
		assert.Equal(t, 3, len(payload.Capabilities))

		time.Sleep(1 * time.Second)

		payload, err = getCapabilities(
			context.Background(), tc.Cachestore(), client, testDomain,
		)
		require.NoError(t, err)
		require.NotNil(t, payload)
		assert.Equal(t, paymail.DefaultBsvAliasVersion, payload.BsvAlias)
		assert.Equal(t, 3, len(payload.Capabilities))
	})
}
