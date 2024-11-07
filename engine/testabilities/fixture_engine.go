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
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

const inMemoryDbConnectionString = "file:spv-wallet-test.db?mode=memory"

type ConfigOpts func(*config.AppConfig)

type EngineFixture interface {
	Engine() (walletEngine EngineWithConfig, cleanup func())
	EngineWithConfiguration(opts ...ConfigOpts) (walletEngine EngineWithConfig, cleanup func())

	// ConfigForTests returns a configuration with default values for tests and with the provided options applied.
	ConfigForTests(opts ...ConfigOpts) *config.AppConfig

	// NewTest creates a new test fixture based on the current one and the provided testing.TB
	// This is useful if you want to start spv-wallet once and then run multiple t.Run with some calls against this one instance.
	NewTest(t testing.TB) EngineFixture

	// BHS creates a new test fixture for Block Header Service (BHS)
	BHS() BlockHeadersServiceFixture
}

type EngineWithConfig struct {
	Config config.AppConfig
	Engine engine.ClientInterface
}

type engineFixture struct {
	config             *config.AppConfig
	engine             engine.ClientInterface
	t                  testing.TB
	logger             zerolog.Logger
	dbConnectionString string
	dbConnection       *sql.DB
	externalTransport  *httpmock.MockTransport
	paymailClient      *paymailmock.PaymailClientMock
}

func Given(t testing.TB) EngineFixture {
	f := &engineFixture{
		t:                 t,
		logger:            tester.Logger(t),
		externalTransport: httpmock.NewMockTransport(),
		// TODO reuse externalTransport in paymailmock
		paymailClient: paymailmock.MockClient(fixtures.PaymailDomainExternal),
	}

	return f
}

func (f *engineFixture) NewTest(t testing.TB) EngineFixture {
	newFixture := *f
	newFixture.t = t
	newFixture.logger = tester.Logger(t)
	return &newFixture
}

func (f *engineFixture) Engine() (walletEngine EngineWithConfig, cleanup func()) {
	return f.EngineWithConfiguration()
}

func (f *engineFixture) EngineWithConfiguration(opts ...ConfigOpts) (walletEngine EngineWithConfig, cleanup func()) {
	f.config = f.ConfigForTests(opts...)
	f.initDbConnection()

	options, err := f.config.ToEngineOptions(f.logger)
	require.NoError(f.t, err)
	options = f.addMockedExternalDependenciesOptions(options)

	f.engine, err = engine.NewClient(context.Background(), options...)
	require.NoError(f.t, err)

	f.initialiseFixtures()

	cleanup = func() {
		err := f.engine.Close(context.Background())
		require.NoError(f.t, err)
		f.externalTransport.Reset()
		httpmock.Reset()
	}

	return EngineWithConfig{Config: *f.config, Engine: f.engine}, cleanup
}

func (f *engineFixture) ConfigForTests(opts ...ConfigOpts) *config.AppConfig {
	configuration := getConfigForTests()

	for _, opt := range opts {
		opt(configuration)
	}

	return configuration
}

// initDbConnection creates a new connection that will be used as connection for engine
func (f *engineFixture) initDbConnection() {
	if f.config.Db.Datastore.Engine != datastore.SQLite {
		panic("Other datastore engines are not supported in tests (yet)")
	}
	// Setting this to give a clue in debugging
	f.dbConnectionString = inMemoryDbConnectionString
	f.config.Db.SQLite.DatabasePath = "already_set_with_existing_connection_config"

	connection, err := sql.Open("sqlite3", f.dbConnectionString)
	require.NoErrorf(f.t, err, "Cannot create sqlite connection")

	f.dbConnection = connection
	f.config.Db.SQLite.ExistingConnection = connection
}

func (f *engineFixture) initialiseFixtures() {
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

func (f *engineFixture) addMockedExternalDependenciesOptions(options []engine.ClientOps) []engine.ClientOps {
	options = append(options, engine.WithHTTPClient(f.httpClientWithMockedTransport()))
	options = append(options, engine.WithPaymailClient(f.paymailClient))
	return options
}

func (f *engineFixture) httpClientWithMockedTransport() *resty.Client {
	client := resty.New()
	client.SetTransport(f.externalTransport)
	return client
}

func getConfigForTests() *config.AppConfig {
	cfg := config.GetDefaultAppConfig()
	cfg.Authentication.RequireSigning = false

	cfg.DebugProfiling = false

	cfg.CustomFeeUnit = &config.FeeUnitConfig{
		Satoshis: 1,
		Bytes:    1000,
	}

	cfg.Paymail.Domains = []string{fixtures.PaymailDomain}

	cfg.Notifications.Enabled = false

	cfg.Db.Datastore.Engine = datastore.SQLite
	cfg.Db.SQLite.DatabasePath = inMemoryDbConnectionString
	cfg.Db.SQLite.TablePrefix = "xapi"
	cfg.Db.SQLite.MaxIdleConnections = 1
	cfg.Db.SQLite.MaxOpenConnections = 1

	return cfg
}