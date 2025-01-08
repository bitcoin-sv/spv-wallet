package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
)

type contactFixture struct {
	engineFixture testengine.EngineFixture
	t             testing.TB
}

type ContactFixture interface {
	Engine() (engine.ClientInterface, func())
	PaymailClient() *paymailmock.PaymailClientMock
}

func given(t testing.TB) ContactFixture {
	return &contactFixture{
		t:             t,
		engineFixture: testengine.Given(t),
	}
}

func (cf *contactFixture) Engine() (engine.ClientInterface, func()) {
	engine, cleanup := cf.engineFixture.Engine()

	return engine.Engine, cleanup
}

func (cf *contactFixture) PaymailClient() *paymailmock.PaymailClientMock {
	return cf.engineFixture.PaymailClient()
}
