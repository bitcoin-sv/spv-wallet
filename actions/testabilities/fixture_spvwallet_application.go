package testabilities

import (
	"context"
	"database/sql"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/server"
	"github.com/go-resty/resty/v2"
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
}

type SPVWalletHttpClientFixture interface {
	// ForAnonymous returns a new http client that is configured without any authentication.
	ForAnonymous() *resty.Client
	// ForAdmin returns a new http client that is configured with the authentication with default admin xpub.
	ForAdmin() *resty.Client
	// ForUser returns a new http client that is configured with the authentication with the xpub of the sender user.
	ForUser() *resty.Client
}

type appFixture struct {
	config             *config.AppConfig
	engine             engine.ClientInterface
	t                  testing.TB
	logger             zerolog.Logger
	transport          testServer
	dbConnectionString string
	dbConnection       *sql.DB
}

func Given(t testing.TB) SPVWalletApplicationFixture {
	f := &appFixture{
		t:      t,
		logger: tester.Logger(t),
		config: getConfigForTests(),
	}

	f.initDbConnection()

	return f
}

func (f *appFixture) NewTest(t testing.TB) SPVWalletApplicationFixture {
	return newOf(*f, t)
}

func newOf(f appFixture, t testing.TB) *appFixture {
	f.t = t
	f.logger = tester.Logger(t)
	return &f
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

	f.engine, err = engine.NewClient(context.Background(), options...)
	require.NoError(f.t, err)

	f.registerUsersFromFixture()

	s := server.NewServer(f.config, f.engine, f.logger)
	f.transport.handlers = s.Handlers()

	return func() {
		err := f.engine.Close(context.Background())
		require.NoError(f.t, err)
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
	c.GetClient().Transport = f.transport
	return c
}

func (f *appFixture) ForAdmin() *resty.Client {
	c := f.ForAnonymous()
	c.SetHeader("x-auth-xpub", config.DefaultAdminXpub)
	return c
}

func (f *appFixture) ForUser() *resty.Client {
	c := f.ForAnonymous()
	c.SetHeader("x-auth-xpub", fixtures.Sender.XPub())
	return c
}

func (f *appFixture) registerUsersFromFixture() {
	opts := f.engine.DefaultModelOptions(engine.WithMetadata("source", "fixture"))

	for _, user := range fixtures.InternalUsers() {
		_, err := f.engine.NewXpub(context.Background(), user.XPub(), opts...)
		require.NoError(f.t, err)

		for _, paymail := range user.Paymails {
			_, err := f.engine.NewPaymailAddress(context.Background(), user.XPub(), paymail, paymail, "", opts...)
			require.NoError(f.t, err)
		}
	}
}

// initDbConnection creates a new connection that will be used as connection for engine
func (f *appFixture) initDbConnection() {
	hex, err := utils.RandomHex(8)
	require.NoErrorf(f.t, err, "cannot generate random hex for sqlite inmemory db name")

	f.dbConnectionString = "file:" + hex + "?mode=memory&cache=shared"

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
