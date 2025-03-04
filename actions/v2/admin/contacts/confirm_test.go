package contacts_test

import (
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"testing"
)

func TestConfirmContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("No side has contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).HasStatus(500).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrGetContact.Code,
				"message": spverrors.ErrGetContact.Message,
			})
	})

	t.Run("Only one side has contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).HasStatus(500).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrGetContact.Code,
				"message": spverrors.ErrGetContact.Message,
			})
	})

	t.Run("Confirm contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).HasStatus(200)
	})

	t.Run("Confirm already confirmed contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).HasStatus(200)

		// and:
		// when:
		res, _ = client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).HasStatus(200)
	})

	t.Run("Confirm contact with user xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.RecipientExternal.DefaultPaymail().String(),
				"paymailB": fixtures.Sender.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

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

}
