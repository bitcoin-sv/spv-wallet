package contacts_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
)

func TestUpdateContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	fullNameForUpdate := "updated full name"

	t.Run("Update contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"fullName": fullNameForUpdate,
			}).
			Put(fmt.Sprintf("/api/v2/admin/contacts/%d", contact.ID))

		// then:
		then.Response(res).
			HasStatus(200).
			WithJSONMatching(`{
				"id": "{{ matchNumber }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"fullName": "{{ .fullName }}",
				"paymail": "{{ .paymail }}",
				"pubKey": "{{ matchHexWithLength 66 }}",
				"status": "{{ .status }}"
			}`, map[string]any{
				"fullName": fullNameForUpdate,
				"paymail":  fixtures.RecipientInternal.DefaultPaymail().String(),
				"status":   contactsmodels.ContactNotConfirmed,
			})
	})

	t.Run("No body in request", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Put(fmt.Sprintf("/api/v2/admin/contacts/%d", contact.ID))

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-bind-body-invalid", "cannot bind request body"))
	})

	t.Run("Update contact with user xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			Put(fmt.Sprintf("/api/v2/admin/contacts/%d", contact.ID))

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONf(apierror.ExpectedJSON("error-unauthorized-xpub-not-an-admin-key", "xpub provided is not an admin key"))
	})

	t.Run("No contact to update", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"fullName": fullNameForUpdate,
			}).
			Put("/api/v2/admin/contacts/99999")

		// then:
		then.Response(res).
			HasStatus(500).
			WithJSONf(apierror.ExpectedJSON("error-contact-updating-status-failed", "updating contact status failed"))
	})
}
