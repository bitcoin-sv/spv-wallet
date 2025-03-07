package contacts_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestDeleteContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("Delete contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/admin/contacts/%d", contact.ID))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Delete already deleted contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/admin/contacts/%d", contact.ID))

		// then:
		then.Response(res).IsOK()

		// when:
		res, _ = client.R().
			Delete(fmt.Sprintf("/api/v2/admin/contacts/%d", contact.ID))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Delete contact with user xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/admin/contacts/%d", contact.ID))

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONf(apierror.ExpectedJSON("error-unauthorized-xpub-not-an-admin-key", "xpub provided is not an admin key"))
	})

	t.Run("No contact to delete", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete("/api/v2/admin/contacts/99999")

		// then:
		then.Response(res).IsOK()
	})
}
