package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

// ActorTestContext provides a simplified test context with actor support
type ActorTestContext struct {
	Given          testabilities.SPVWalletApplicationFixture
	baseAssertions testabilities.SPVWalletApplicationAssertions
	t              testing.TB

	Alice   *User
	Bob     *User
	Charlie *User
}

// User represents either an internal or external wallet user
type User struct {
	fixtures.User
	app testabilities.SPVWalletApplicationFixture
	t   testing.TB
}

// ActorAssertions wraps assertions for a specific actor
type ActorAssertions struct {
	userAssertions testabilities.SPVWalletAppUserAssertions
}

// NewActorTests creates a new test context with actor support
func NewActorTests(t testing.TB) *ActorTestContext {
	given, then := testabilities.New(t)

	ctx := &ActorTestContext{
		Given:          given,
		baseAssertions: then,
		t:              t,
	}

	ctx.Alice = &User{
		User: fixtures.Sender,
		app:  given,
		t:    t,
	}

	ctx.Bob = &User{
		User: fixtures.RecipientInternal,
		app:  given,
		t:    t,
	}

	ctx.Charlie = &User{
		User: fixtures.RecipientExternal,
		app:  given,
		t:    t,
	}

	return ctx
}
