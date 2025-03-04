package contacts_test

import (
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"strconv"
	"testing"
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
			Put("/api/v2/admin/contacts/" + strconv.Itoa(int(contact.ID)))

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
				"status":   "unconfirmed",
			})
	})

	t.Run("No body in request", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Put("/api/v2/admin/contacts/" + strconv.Itoa(int(contact.ID)))

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrCannotBindRequest.Code,
				"message": spverrors.ErrCannotBindRequest.Message,
			})
	})

	t.Run("Update contact with user xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			Put("/api/v2/admin/contacts/" + strconv.Itoa(int(contact.ID)))

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrNotAnAdminKey.Code,
				"message": spverrors.ErrNotAnAdminKey.Message,
			})
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
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrUpdateContactStatus.Code,
				"message": spverrors.ErrUpdateContactStatus.Message,
			})
	})
}
