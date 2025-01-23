package testabilities

import (
	"context"
	"errors"
	"testing"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/testabilities/testmode"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

const inMemoryDbConnectionString = "file:spv-wallet-test.db?mode=memory"
const fileDbConnectionString = "file:spv-wallet-test.db"

type EngineFixture interface {
	Engine() (walletEngine EngineWithConfig, cleanup func())
	EngineWithConfiguration(opts ...ConfigOpts) (walletEngine EngineWithConfig, cleanup func())
	PaymailClient() *paymailmock.PaymailClientMock

	// ConfigForTests returns a configuration with default values for tests and with the provided options applied.
	ConfigForTests(opts ...ConfigOpts) *config.AppConfig

	// NewTest creates a new test fixture based on the current one and the provided testing.TB
	// This is useful if you want to start spv-wallet once and then run multiple t.Run with some calls against this one instance.
	NewTest(t testing.TB) EngineFixture

	// BHS creates a new test fixture for Block Header Service (BHS)
	BHS() BlockHeadersServiceFixture

	// ARC creates a new test fixture for ARC
	ARC() ARCFixture
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

func (f *engineFixture) PaymailClient() *paymailmock.PaymailClientMock {
	return f.paymailClient
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
	f.prepareDBConfigForTests()

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

	return EngineWithConfig{
		Config: *f.config,
		Engine: f.engine,
	}, cleanup
}

func (f *engineFixture) ConfigForTests(opts ...ConfigOpts) *config.AppConfig {
	configuration := getConfigForTests()

	for _, opt := range opts {
		opt(configuration)
	}

	return configuration
}

// prepareDBConfigForTests creates a new connection that will be used as connection for engine
func (f *engineFixture) prepareDBConfigForTests() {
	require.Equal(f.t, datastore.SQLite, f.config.Db.Datastore.Engine, "Other datastore engines are not supported in tests (yet)")

	// It is a workaround for development purpose to check the code with postgres instance.
	if ok, dbName := testmode.CheckPostgresMode(); ok {
		f.config.Db.Datastore.Engine = datastore.PostgreSQL
		f.config.Db.SQL.User = "postgres"
		f.config.Db.SQL.Password = "postgres"
		f.config.Db.SQL.Name = dbName
		f.config.Db.SQL.Host = "localhost"
		return
	}

	// It is a workaround for development purpose to check what is the db state after running a tests.
	if testmode.CheckFileSQLiteMode() {
		f.dbConnectionString = fileDbConnectionString
	} else {
		f.dbConnectionString = inMemoryDbConnectionString
	}
	f.config.Db.SQLite.Shared = false
	f.config.Db.SQLite.MaxIdleConnections = 1
	f.config.Db.SQLite.MaxOpenConnections = 1
	f.config.Db.SQLite.DatabasePath = f.dbConnectionString
}

func (f *engineFixture) initialiseFixtures() {
	opts := f.engine.DefaultModelOptions(engine.WithMetadata("source", "fixture"))

	for _, user := range fixtures.InternalUsers() {
		_, err := f.engine.NewXpub(context.Background(), user.XPub(), opts...)
		if !errors.Is(err, spverrors.ErrXPubAlreadyExists) {
			require.NoError(f.t, err)
		}

		for _, address := range user.Paymails {
			_, err := f.engine.NewPaymailAddress(context.Background(), user.XPub(), address, address, "", opts...)
			if !errors.Is(err, spverrors.ErrPaymailAlreadyExists) {
				require.NoError(f.t, err)
			}
		}

		if f.config.ExperimentalFeatures.NewTransactionFlowEnabled {
			pubKeyHex := user.PublicKey().ToDERHex()
			createdUser, err := f.engine.UsersService().Create(context.Background(), &usersmodels.NewUser{
				PublicKey: pubKeyHex,
			})
			require.NoError(f.t, err)
			for _, address := range user.Paymails {
				alias, domain, _ := paymail.SanitizePaymail(address)
				_, err = f.engine.PaymailsService().Create(context.Background(), &paymailsmodels.NewPaymail{
					Alias:  alias,
					Domain: domain,

					PublicName: address,
					Avatar:     "",
					UserID:     createdUser.ID,
				})
			}
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
