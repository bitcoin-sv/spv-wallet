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

func TestCreateContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("Create contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"creatorPaymail": fixtures.Sender.DefaultPaymail().String(),
				"fullName":       fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Post(fmt.Sprintf("/api/v2/admin/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

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
				"fullName": fixtures.RecipientInternal.DefaultPaymail().PublicName(),
				"paymail":  fixtures.RecipientInternal.DefaultPaymail().String(),
				"status":   contactsmodels.ContactNotConfirmed,
			})
	})

	t.Run("Create contact with user xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"creatorPaymail": fixtures.Sender.DefaultPaymail().String(),
				"fullName":       fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Post(fmt.Sprintf("/api/v2/admin/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONf(apierror.ExpectedJSON("error-unauthorized-xpub-not-an-admin-key", "xpub provided is not an admin key"))
	})

	t.Run("Create contact with unknown creator paymail", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"creatorPaymail": "unknown-paymail@exmaple.com",
				"fullName":       fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Post(fmt.Sprintf("/api/v2/admin/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		fmt.Println(res.String())

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONf(apierror.ExpectedJSON("error-paymail-not-found", "paymail not found"))
	})

	t.Run("Create contact with unknown requester paymail", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"creatorPaymail": fixtures.Sender.DefaultPaymail().String(),
				"fullName":       fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Post("/api/v2/admin/contacts/unknown-paymail@exmaple.com")

		fmt.Println(res.String())

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-contact-getting-pki-failed", "getting PKI for contact failed"))
	})

	t.Run("Create contact without creator paymail", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"fullName": fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Post("/api/v2/admin/contacts/unknown-paymail@exmaple.com")

		fmt.Println(res.String())

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-contact-creator-paymail-missing", "missing creator paymail in contact"))
	})

	t.Run("Create contact without contact full name", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"creatorPaymail": fixtures.Sender.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/unknown-paymail@exmaple.com")

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-contact-full-name-missing", "missing full name in contact"))
	})
}
