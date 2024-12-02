package users_test

import (
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"testing"
)

func TestCurrentUserPaymails(t *testing.T) {
	t.Run("return paymails info for user (single paymail)", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().Get("/api/v1/users/current/paymails")

		// then:
		then.Response(res).
			IsOK()
	})

	t.Run("return paymails info for user (multiple paymails)", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForGivenUser(fixtures.UserWithMorePaymails)

		// when:
		res, _ := client.R().Get("/api/v1/users/current/paymails")

		// then:
		then.Response(res).
			IsOK()
	})

	t.Run("try to return paymails info for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v1/users/current/paymails")

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("return xpub info for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get("/api/v1/users/current/paymails")

		// then:
		then.Response(res).IsUnauthorized()
	})
}
