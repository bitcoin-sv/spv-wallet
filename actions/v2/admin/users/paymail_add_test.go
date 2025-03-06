package users_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestAddPaymail(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	user := fixtures.Sender

	// and:
	secondPaymail := fixtures.Paymail("sender_second@" + fixtures.PaymailDomain)

	// and:
	thirdPaymail := fixtures.Paymail("third" + "@" + fixtures.PaymailDomain)

	// and:
	avatarURL := "https://address-to-avatar.com"

	t.Run("Add a paymail to a user as admin using whole paymail address", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"address":    secondPaymail,
				"publicName": secondPaymail.PublicName(),
				"avatarURL":  avatarURL,
			}).
			SetPathParam("id", user.ID()).
			Post("/api/v2/admin/users/{id}/paymails")

		// then:
		then.Response(res).
			HasStatus(201).
			WithJSONMatching(`{
			  "alias": "{{ .alias }}",
			  "avatar": "{{ .avatar }}",
			  "domain": "example.com",
			  "id": "{{ matchNumber }}",
			  "paymail": "{{ .paymail }}",
			  "publicName": "{{ .publicName }}"
			}`, map[string]any{
				"paymail":    secondPaymail,
				"publicName": secondPaymail.PublicName(),
				"alias":      secondPaymail.Alias(),
				"avatar":     avatarURL,
			})
	})

	t.Run("Try to add a paymail to a user as admin using whole paymail address and wrong url avatar", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()
		wrongAvatarURL := "/User/path/to/avatar"

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"address":    secondPaymail,
				"publicName": secondPaymail.PublicName(),
				"avatarURL":  wrongAvatarURL,
			}).
			SetPathParam("id", user.ID()).
			Post("/api/v2/admin/users/{id}/paymails")

		// then:
		then.Response(res).
			WithProblemDetails(422, "invalid_avatar_url", "Invalid avatar URL")
	})

	t.Run("Add a paymail to a user as admin using alias and domain as address", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"alias":  thirdPaymail.Alias(),
				"domain": thirdPaymail.Domain(),
			}).
			SetPathParam("id", user.ID()).
			Post("/api/v2/admin/users/{id}/paymails")

		// then:
		then.Response(res).
			HasStatus(201).
			WithJSONMatching(`{
			  "alias": "{{ .alias }}",
			  "avatar": "",
			  "domain": "example.com",
			  "id": "{{ matchNumber }}",
			  "paymail": "{{ .paymail }}",
			  "publicName": "{{ .alias }}"
			}`, map[string]any{
				"paymail": thirdPaymail,
				"alias":   thirdPaymail.Alias(),
			})
	})

	t.Run("Get user with second paymail by id as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetPathParam("id", user.ID()).
			Get("/api/v2/admin/users/{id}")

		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"id": "{{ matchAddress }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"publicKey": "{{ .publicKey }}",
				"paymails": [
					{
						"alias": "{{ .defaultAlias }}",
						"avatar": "",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .defaultPaymail }}",
						"publicName": "{{ .defaultPublicName }}"
					},
					{
						"alias": "{{ .secondAlias }}",
						"avatar": "{{ .secondAvatar }}",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .secondPaymail }}",
						"publicName": "{{ .secondPublicName }}"
					},
					{
						"alias": "{{ .thirdAlias }}",
						"avatar": "",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .thirdPaymail }}",
						"publicName": "{{ .thirdAlias }}"
					}
				]
			}`, map[string]any{
				"publicKey": user.PublicKey().ToDERHex(),

				"defaultPaymail":    user.DefaultPaymail(),
				"defaultPublicName": user.DefaultPaymail().PublicName(),
				"defaultAlias":      user.DefaultPaymail().Alias(),

				"secondPaymail":    secondPaymail,
				"secondPublicName": secondPaymail.PublicName(),
				"secondAlias":      secondPaymail.Alias(),
				"secondAvatar":     avatarURL,

				"thirdPaymail": thirdPaymail,
				"thirdAlias":   thirdPaymail.Alias(),
			})
	})
}

func TestAddPaymailWithWrongDomain(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	alias := "user"
	unsupportedDomain := "unsupported.com"
	unsupportedPaymail := alias + "@" + unsupportedDomain
	publicName := "User"

	// and:
	client := given.HttpClient().ForAdmin()

	t.Run("Try to add using whole paymail as address", func(t *testing.T) {
		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"address":    unsupportedPaymail,
				"publicName": publicName,
				"avatarURL":  "",
			}).
			SetPathParam("id", fixtures.Sender.ID()).
			Post("/api/v2/admin/users/{id}/paymails")

		// then:
		then.Response(res).
			WithProblemDetails(400, "unsupported_domain", "Unsupported domain")
	})

	t.Run("Try to add using alias and domain as address", func(t *testing.T) {
		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"alias":      alias,
				"domain":     unsupportedDomain,
				"publicName": publicName,
				"avatar":     "",
			}).
			SetPathParam("id", fixtures.Sender.ID()).
			Post("/api/v2/admin/users/{id}/paymails")

		// then:
		then.Response(res).
			WithProblemDetails(400, "unsupported_domain", "Unsupported domain")
	})

}

func TestAddPaymailWithBothPaymailAndAliasDomainPair(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	client := given.HttpClient().ForAdmin()

	t.Run("Add using consistent fields", func(t *testing.T) {
		// given:
		alias := "user"
		domain := fixtures.PaymailDomain
		paymail := alias + "@" + domain
		publicName := "User"

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"address":    paymail,
				"alias":      alias,
				"domain":     domain,
				"publicName": publicName,
				"avatarURL":  "",
			}).
			SetPathParam("id", fixtures.Sender.ID()).
			Post("/api/v2/admin/users/{id}/paymails")

		// then:
		then.Response(res).
			HasStatus(201)
	})

	t.Run("Try to add with inconsistent paymail and alias-domain pair", func(t *testing.T) {
		// given:
		alias := "user"
		domain := fixtures.PaymailDomain
		paymail := "other" + "@" + domain
		publicName := "User"

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"alias":      alias,
				"domain":     domain,
				"address":    paymail,
				"publicName": publicName,
				"avatarURL":  "",
			}).
			SetPathParam("id", fixtures.Sender.ID()).
			Post("/api/v2/admin/users/{id}/paymails")

		// then:
		then.Response(res).
			HasStatus(400).
			WithProblemDetails(400, "inconsistent_alias_domain_and_address", "Inconsistent alias@domain and address fields")
	})
}
