package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
)

type contactFixture struct {
	engineFixture testengine.EngineFixture
	t             testing.TB
}

type ContactFixture interface {
	Engine() (engine.ClientInterface, func())
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
