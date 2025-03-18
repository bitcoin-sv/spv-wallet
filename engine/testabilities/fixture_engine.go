package testabilities

import (
	"context"
	"errors"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/datastore"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/testabilities/testmode"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/initializer"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const inMemoryDbConnectionString = "file:spv-wallet-test.db?mode=memory"
const fileDbConnectionString = "file:spv-wallet-test.db"

const CallbackTestToken = "arc-test-token"

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

	// Faucet creates a new test fixture for Faucet
	Faucet(user fixtures.User) FaucetFixture

	// User creates a new test fixture for Contacts
	User(user fixtures.User) UserFixture

	// Tx creates a new mocked transaction builder
	Tx() txtestability.TransactionSpec
}

// FaucetFixture is a test fixture for the faucet service
type FaucetFixture interface {
	TopUp(satoshis bsv.Satoshis) txtestability.TransactionSpec
	StoreData(data string) (txtestability.TransactionSpec, string)
}

// UserFixture is a test fixture for the user management
type UserFixture interface {
	HasContactTo(userB fixtures.User) *contactsmodels.Contact
	HasConfirmedContactTo(userB fixtures.User) *contactsmodels.Contact
	HasRejectedContactTo(userB fixtures.User) *contactsmodels.Contact
	HasAwaitingContactTo(userB fixtures.User) *contactsmodels.Contact
}

type EngineWithConfig struct {
	Config config.AppConfig
	Engine engine.ClientInterface
}

type engineFixture struct {
	config                       *config.AppConfig
	engine                       engine.ClientInterface
	t                            testing.TB
	logger                       zerolog.Logger
	dbConnectionString           string
	externalTransport            *httpmock.MockTransport
	paymailClient                *paymailmock.PaymailClientMock
	txFixture                    txtestability.TransactionsFixtures
	externalTransportWithSniffer *tester.HTTPSniffer
}

func Given(t testing.TB) EngineFixture {
	externalTransport := httpmock.NewMockTransport()
	f := &engineFixture{
		t:                            t,
		logger:                       tester.Logger(t),
		externalTransport:            externalTransport,
		paymailClient:                paymailmock.MockClient(externalTransport, fixtures.PaymailDomainExternal),
		txFixture:                    txtestability.Given(t),
		externalTransportWithSniffer: tester.NewHTTPSniffer(externalTransport),
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
	newFixture.txFixture = txtestability.Given(t)
	return &newFixture
}

func (f *engineFixture) Engine() (walletEngine EngineWithConfig, cleanup func()) {
	return f.EngineWithConfiguration()
}

func (f *engineFixture) EngineWithConfiguration(opts ...ConfigOpts) (walletEngine EngineWithConfig, cleanup func()) {
	f.config = f.ConfigForTests(opts...)
	f.prepareDBConfigForTests()

	options, err := initializer.ToEngineOptions(f.config, f.logger)
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

func (f *engineFixture) Faucet(user fixtures.User) FaucetFixture {
	return &faucetFixture{
		engine:    f.engine,
		user:      user,
		t:         f.t,
		assert:    assert.New(f.t),
		require:   require.New(f.t),
		arc:       f.ARC(),
		bhs:       f.BHS(),
		txFixture: f.txFixture,
	}
}

func (f *engineFixture) User(user fixtures.User) UserFixture {
	return &contactFixture{
		engine:        f.engine,
		user:          user,
		t:             f.t,
		assert:        assert.New(f.t),
		require:       require.New(f.t),
		paymailClient: f.paymailClient,
	}
}

func (f *engineFixture) Tx() txtestability.TransactionSpec {
	return f.txFixture.Tx()
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
		if f.config.ExperimentalFeatures.V2 {
			f.createV2User(user, opts)
		} else {
			f.createV1User(user, opts)
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
	client.SetTransport(f.externalTransportWithSniffer)
	return client
}

func (f *engineFixture) createV1User(user fixtures.User, opts []engine.ModelOps) {
	_, err := f.engine.NewXpub(context.Background(), user.XPub(), opts...)
	if !errors.Is(err, spverrors.ErrXPubAlreadyExists) {
		require.NoError(f.t, err)
	}

	for _, pm := range user.Paymails {
		_, err := f.engine.NewPaymailAddress(context.Background(), user.XPub(), pm.Address(), pm.PublicName(), "", opts...)
		if !errors.Is(err, spverrors.ErrPaymailAlreadyExists) {
			require.NoError(f.t, err)
		}
	}
}

func (f *engineFixture) createV2User(user fixtures.User, opts []engine.ModelOps) {
	exists, err := f.engine.UsersService().Exists(context.Background(), user.ID())
	require.NoError(f.t, err)
	if exists {
		return
	}

	pubKeyHex := user.PublicKey().ToDERHex()

	createdUser, err := f.engine.UsersService().Create(context.Background(), &usersmodels.NewUser{
		PublicKey: pubKeyHex,
	})
	require.NoError(f.t, err)
	for _, pm := range user.Paymails {
		_, err = f.engine.PaymailsService().Create(context.Background(), &paymailsmodels.NewPaymail{
			Alias:  pm.Alias(),
			Domain: pm.Domain(),

			PublicName: pm.PublicName(),
			Avatar:     "",
			UserID:     createdUser.ID,
		})
	}
	require.NoError(f.t, err)
}

func getConfigForTests() *config.AppConfig {
	cfg := config.GetDefaultAppConfig()
	cfg.Authentication.RequireSigning = false

	cfg.DebugProfiling = false

	cfg.CustomFeeUnit = &config.FeeUnitConfig{
		Satoshis: 1,
		Bytes:    1000,
	}

	cfg.ARC.Callback.Enabled = true
	cfg.ARC.Callback.Host = "https://" + fixtures.PaymailDomain
	cfg.ARC.Callback.Token = CallbackTestToken

	cfg.Paymail.Domains = []string{fixtures.PaymailDomain}

	cfg.Notifications.Enabled = false

	cfg.Db.Datastore.Engine = datastore.SQLite
	cfg.Db.SQLite.DatabasePath = inMemoryDbConnectionString
	cfg.Db.SQLite.TablePrefix = "xapi"
	cfg.Db.SQLite.MaxIdleConnections = 1
	cfg.Db.SQLite.MaxOpenConnections = 1

	return cfg
}
