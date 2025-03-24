package contacts_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestRejectContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("Reject contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/admin/invitations/%d", contact.ID))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Reject already rejected contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/admin/invitations/%d", contact.ID))

		// then:
		then.Response(res).IsOK()

		// when:
		res, _ = client.R().
			Delete(fmt.Sprintf("/api/v2/admin/invitations/%d", contact.ID))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Reject contact with user xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/admin/invitations/%d", contact.ID))

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONf(apierror.ExpectedJSON("error-unauthorized-xpub-not-an-admin-key", "xpub provided is not an admin key"))
	})

	t.Run("No contact to reject", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete("/api/v2/admin/invitations/99999")

		// then:
		then.Response(res).
			HasStatus(500).
			WithJSONf(apierror.ExpectedJSON("error-contact-updating-status-failed", "updating contact status failed"))
	})
}
