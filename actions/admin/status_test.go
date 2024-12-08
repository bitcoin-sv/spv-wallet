package admin_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
)

func TestGETAdminStatus(t *testing.T) {
	t.Run("return success if authenticated as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)

		cleanup := given.StartedSPVWallet()
		defer cleanup()

		client := given.HttpClient().ForAdmin()
		// when:
		res, _ := client.R().Get("/api/v1/admin/status")

		// then:
		then.Response(res).IsOK()
	})

	t.Run("return unauthorized if not authenticated as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)

		cleanup := given.StartedSPVWallet()
		defer cleanup()

		client := given.HttpClient().ForUser()
		// when:
		res, _ := client.R().Get("/api/v1/admin/status")

		// then:
		then.Response(res).IsUnauthorizedForUser()
	})

	t.Run("return unauthorized if not authenticated at all", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)

		cleanup := given.StartedSPVWallet()
		defer cleanup()

		client := given.HttpClient().ForAnonymous()
		// when:
		res, _ := client.R().Get("/api/v1/admin/status")

		// then:
		then.Response(res).IsUnauthorized()
	})
}
