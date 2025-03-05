package contacts_test

import (
	"fmt"

	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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
			Post("/api/v2/admin/contacts/" + fixtures.RecipientInternal.DefaultPaymail().String())

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
			Post("/api/v2/admin/contacts/" + fixtures.RecipientInternal.DefaultPaymail().String())

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
			Post("/api/v2/admin/contacts/" + fixtures.RecipientInternal.DefaultPaymail().String())

		fmt.Println(res.String())

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrCouldNotFindPaymail.Code,
				"message": spverrors.ErrCouldNotFindPaymail.Message,
			})
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
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrGettingPKIFailed.Code,
				"message": spverrors.ErrGettingPKIFailed.Message,
			})
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
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrMissingContactCreatorPaymail.Code,
				"message": spverrors.ErrMissingContactCreatorPaymail.Message,
			})
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
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrMissingContactFullName.Code,
				"message": spverrors.ErrMissingContactFullName.Message,
			})
	})
}
