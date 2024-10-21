package testabilities

import (
	"context"
	"database/sql"
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/bitcoin-sv/spv-wallet/server"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type ConfigOpts func(*config.AppConfig)

type SPVWalletApplicationFixture interface {
	StartedSPVWallet() (cleanup func())
	StartedSPVWalletWithConfiguration(opts ...ConfigOpts) (cleanup func())

	// HttpClient returns a new http client that can be used to make requests to the spv-wallet server.
	// It is also failing tests if there are network or invalid request configuration errors,
	// so the tests can focus on the server response.
	HttpClient() SPVWalletHttpClientFixture

	// NewTest creates a new test fixture based on the current one and the provided testing.TB
	// This is useful if you want to start spv-wallet once and then run multiple t.Run with some calls against this one instance.
	NewTest(t testing.TB) SPVWalletApplicationFixture

	// BHS creates a new test fixture for Block Header Service (BHS)
	BHS() BlockHeadersServiceFixture
}

type BlockHeadersServiceFixture interface {
	// WillRespondForMerkleRoots returns a http response for get merkleroots endpoint with
	// provided httpCode and response
	WillRespondForMerkleRoots(httpCode int, response string)
}

type SPVWalletHttpClientFixture interface {
	// ForAnonymous returns a new http client that is configured without any authentication.
	ForAnonymous() *resty.Client
	// ForAdmin returns a new http client that is configured with the authentication with default admin xpub.
	ForAdmin() *resty.Client
	// ForUser returns a new http client that is configured with the authentication with the xpub of the sender user.
	ForUser() *resty.Client
	// ForGivenUser returns a new http client that is configured with the authentication with the xpub of the given user.
	ForGivenUser(user fixtures.User) *resty.Client
}

type appFixture struct {
	config             *config.AppConfig
	engine             engine.ClientInterface
	t                  testing.TB
	logger             zerolog.Logger
	server             testServer
	dbConnectionString string
	dbConnection       *sql.DB
	externalTransport  *httpmock.MockTransport
	paymailClient      *paymailmock.PaymailClientMock
}

func Given(t testing.TB) SPVWalletApplicationFixture {
	f := &appFixture{
		t:                 t,
		logger:            tester.Logger(t),
		config:            getConfigForTests(),
		externalTransport: httpmock.NewMockTransport(),
		// TODO reuse externalTransport in paymailmock
		paymailClient: paymailmock.MockClient(fixtures.PaymailDomainExternal),
	}

	f.initDbConnection()

	return f
}

func (f *appFixture) NewTest(t testing.TB) SPVWalletApplicationFixture {
	newFixture := *f
	newFixture.t = t
	newFixture.logger = tester.Logger(t)
	return &newFixture
}

func (f *appFixture) StartedSPVWallet() (cleanup func()) {
	return f.StartedSPVWalletWithConfiguration()
}

func (f *appFixture) StartedSPVWalletWithConfiguration(opts ...ConfigOpts) (cleanup func()) {
	for _, opt := range opts {
		opt(f.config)
	}

	options, err := f.config.ToEngineOptions(f.logger)
	require.NoError(f.t, err)
	options = f.addMockedExternalDependenciesOptions(options)

	f.engine, err = engine.NewClient(context.Background(), options...)
	require.NoError(f.t, err)

	f.initialiseFixtures()

	s := server.NewServer(f.config, f.engine, f.logger)
	f.server.handlers = s.Handlers()

	return func() {
		err := f.engine.Close(context.Background())
		require.NoError(f.t, err)
		f.externalTransport.Reset()
		httpmock.Reset()
	}
}

func (f *appFixture) HttpClient() SPVWalletHttpClientFixture {
	return f
}

func (f *appFixture) ForAnonymous() *resty.Client {
	c := resty.New()
	c.OnError(func(_ *resty.Request, err error) {
		f.t.Fatalf("HTTP request end up with unexpected error: %v", err)
	})
	c.GetClient().Transport = f.server
	return c
}

func (f *appFixture) ForAdmin() *resty.Client {
	c := f.ForAnonymous()
	c.SetHeader("x-auth-xpub", config.DefaultAdminXpub)
	return c
}

func (f *appFixture) ForUser() *resty.Client {
	return f.ForGivenUser(fixtures.Sender)
}

func (f *appFixture) ForGivenUser(user fixtures.User) *resty.Client {
	c := f.ForAnonymous()
	c.SetHeader("x-auth-xpub", user.XPub())
	return c
}

func (f *appFixture) BHS() BlockHeadersServiceFixture {
	return f
}

func (f *appFixture) WillRespondForMerkleRoots(httpCode int, response string) {
	responder := func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(httpCode, response)
		res.Header.Set("Content-Type", "application/json")

		return res, nil
	}

	f.externalTransport.RegisterResponder("GET", "http://localhost:8080/api/v1/chain/merkleroot", responder)
}

func (f *appFixture) mockBHSGetMerkleRoots() {
	responder := func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("Authorization") != "Bearer "+f.config.BHS.AuthToken {
			return httpmock.NewStringResponse(http.StatusUnauthorized, ""), nil
		}
		lastEvaluatedKey := req.URL.Query().Get("lastEvaluatedKey")
		merkleRootsRes, err := simulateBHSMerkleRootsAPI(lastEvaluatedKey)
		require.NoError(f.t, err)

		res := httpmock.NewStringResponse(200, merkleRootsRes)
		res.Header.Set("Content-Type", "application/json")

		return res, nil
	}

	f.externalTransport.RegisterResponder("GET", "http://localhost:8080/api/v1/chain/merkleroot", responder)
}

func (f *appFixture) initialiseFixtures() {
	opts := f.engine.DefaultModelOptions(engine.WithMetadata("source", "fixture"))

	for _, user := range fixtures.InternalUsers() {
		_, err := f.engine.NewXpub(context.Background(), user.XPub(), opts...)
		require.NoError(f.t, err)

		for _, paymail := range user.Paymails {
			_, err := f.engine.NewPaymailAddress(context.Background(), user.XPub(), paymail, paymail, "", opts...)
			require.NoError(f.t, err)
		}
	}

	f.paymailClient.WillRespondWithP2PCapabilities()
	f.mockBHSGetMerkleRoots()
}

// initDbConnection creates a new connection that will be used as connection for engine
func (f *appFixture) initDbConnection() {
	f.dbConnectionString = "file:spv-wallet-test.db?mode=memory"

	connection, err := sql.Open("sqlite3", f.dbConnectionString)
	require.NoErrorf(f.t, err, "Cannot create sqlite connection")

	f.dbConnection = connection

	f.config.Db.Datastore.Engine = datastore.SQLite
	f.config.Db.SQLite.DatabasePath = ""
	f.config.Db.SQLite.TablePrefix = "xapi"
	f.config.Db.SQLite.MaxIdleConnections = 1
	f.config.Db.SQLite.MaxOpenConnections = 1
	f.config.Db.SQLite.ExistingConnection = f.dbConnection
}

func (f *appFixture) addMockedExternalDependenciesOptions(options []engine.ClientOps) []engine.ClientOps {
	options = append(options, engine.WithHTTPClient(f.httpClientWithMockedTransport()))
	options = append(options, engine.WithPaymailClient(f.paymailClient))
	return options
}

func (f *appFixture) httpClientWithMockedTransport() *resty.Client {
	client := resty.New()
	client.SetTransport(f.externalTransport)
	return client
}

func getConfigForTests() *config.AppConfig {
	cfg := config.GetDefaultAppConfig()
	cfg.Authentication.RequireSigning = false

	cfg.DebugProfiling = false

	cfg.ARC.UseFeeQuotes = false
	cfg.ARC.FeeUnit = &config.FeeUnitConfig{
		Satoshis: 1,
		Bytes:    1000,
	}

	cfg.Paymail.Domains = []string{fixtures.PaymailDomain}

	cfg.Notifications.Enabled = false

	return cfg
}
