package contacts_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestRemoveContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("Remove contact in unconfirmed status", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Remove contact in confirmed status", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasConfirmedContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Remove contact in awaiting status", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasAwaitingContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Remove contact in rejected status", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasRejectedContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Remove already removed contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasConfirmedContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()

		// and:
		// when:
		res, _ = client.R().
			Delete(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Remove contact with admin xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONf(apierror.ExpectedJSON("error-admin-auth-on-user-endpoint", "cannot call user's endpoints with admin authorization"))
	})

	t.Run("No contact to remove", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientExternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()
	})
}
