package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/config"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testpaymail "github.com/bitcoin-sv/spv-wallet/engine/paymail/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/server"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type SPVWalletApplicationFixture interface {
	StartedSPVWallet() (cleanup func())
	StartedSPVWalletWithConfiguration(opts ...testengine.ConfigOpts) (cleanup func())

	// HttpClient returns a new http client that can be used to make requests to the spv-wallet server.
	// It is also failing tests if there are network or invalid request configuration errors,
	// so the tests can focus on the server response.
	HttpClient() SPVWalletHttpClientFixture

	// NewTest creates a new test fixture based on the current one and the provided testing.TB
	// This is useful if you want to start spv-wallet once and then run multiple t.Run with some calls against this one instance.
	NewTest(t testing.TB) SPVWalletApplicationFixture

	// BHS creates a new test fixture for Block Header Service (BHS)
	BHS() BlockHeadersServiceFixture

	// ARC creates a new test fixture for ARC
	ARC() ARCFixture

	Paymail() testpaymail.PaymailClientFixture

	Faucet(user fixtures.User) testengine.FaucetFixture

	User(user fixtures.User) testengine.UserFixture

	EngineFixture() testengine.EngineFixture

	// Tx creates a new mocked transaction builder
	Tx() txtestability.TransactionSpec

	Config() *config.AppConfig
}

type BlockHeadersServiceFixture interface {
	// WillRespondForMerkleRoots returns a http response for get merkleroots endpoint with
	// provided httpCode and response
	WillRespondForMerkleRoots(httpCode int, response string)

	// WillRespondForMerkleRootsVerify returns a MerkleRootsConfirmations response for get merkleroot/verify endpoint with
	// provided httpCode
	WillRespondForMerkleRootsVerify(httpCode int, response *chainmodels.MerkleRootsConfirmations)
}

type ARCFixture interface {
	// WillRespondForBroadcast returns a http response for a broadcast request.
	WillRespondForBroadcast(httpCode int, info *chainmodels.TXInfo)

	// WillRespondForBroadcastWithSeenOnNetwork is a default ARC behavior for broadcasting (happy path).
	WillRespondForBroadcastWithSeenOnNetwork(txID string)
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
	engineWithConfig testengine.EngineWithConfig
	engineFixture    testengine.EngineFixture
	t                testing.TB
	logger           zerolog.Logger
	server           testServer
}

func Given(t testing.TB) SPVWalletApplicationFixture {
	f := &appFixture{
		t:             t,
		engineFixture: testengine.Given(t),
		logger:        tester.Logger(t),
	}
	return f
}

func (f *appFixture) NewTest(t testing.TB) SPVWalletApplicationFixture {
	newFixture := *f
	newFixture.t = t
	newFixture.logger = tester.Logger(t)
	newFixture.engineFixture = f.engineFixture.NewTest(t)

	newFixture.engineFixture.PaymailClient().RedirectTransportIfDomain(fixtures.PaymailDomain, f.server)
	return &newFixture
}

func (f *appFixture) StartedSPVWallet() (cleanup func()) {
	return f.StartedSPVWalletWithConfiguration()
}

func (f *appFixture) StartedSPVWalletWithConfiguration(opts ...testengine.ConfigOpts) (cleanup func()) {
	engineWithConfig, cleanup := f.engineFixture.EngineWithConfiguration(opts...)

	s := server.NewServer(&engineWithConfig.Config, engineWithConfig.Engine, f.logger)
	f.server.handlers = s.Handlers()

	f.engineWithConfig = engineWithConfig

	return cleanup
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
	return f.engineFixture.BHS()
}

func (f *appFixture) ARC() ARCFixture {
	return f.engineFixture.ARC()
}

func (f *appFixture) Paymail() testpaymail.PaymailClientFixture {
	return testpaymail.GivenWithMockClient(f.t, f.engineFixture.PaymailClient())
}

func (f *appFixture) Faucet(user fixtures.User) testengine.FaucetFixture {
	return f.engineFixture.Faucet(user)
}

func (f *appFixture) User(user fixtures.User) testengine.UserFixture {
	return f.engineFixture.User(user)
}

func (f *appFixture) Tx() txtestability.TransactionSpec {
	return f.engineFixture.Tx()
}

func (f *appFixture) EngineFixture() testengine.EngineFixture {
	return f.engineFixture
}

func (f *appFixture) Config() *config.AppConfig {
	return &f.engineWithConfig.Config
}
