package paymails_test

import (
	"strings"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestCurrentUserPaymails(t *testing.T) {
	t.Run("return paymails info for user (single paymail)", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().Get("/api/v1/paymails")

		// then:
		then.Response(res).
			IsOK().
			WithJSONTemplate(`{
			 "content": [
				{
				  "address": "{{.Address}}",
				  "alias": "{{.Alias}}",
				  "avatar": "/.*/",
				  "createdAt": "/.*/",
				  "deletedAt": null,
				  "domain": "{{.Domain}}",
				  "id": "/^[a-zA-Z0-9]{64}$/",
				  "metadata": "*",
				  "publicName": "{{.PublicName}}",
				  "updatedAt": "/.*/",
				  "xpubId": "{{.XPubID}}"
				}
			 ],
			 "page": {
				"number": 1,
				"size": 50,
				"totalElements": 1,
				"totalPages": 1
			 }
			}`, map[string]any{
				"Address":    strings.ToLower(fixtures.Sender.Paymails[0]),
				"PublicName": fixtures.Sender.Paymails[0],
				"Alias":      getAliasFromPaymail(t, fixtures.Sender.Paymails[0]),
				"XPubID":     fixtures.Sender.XPubID(),
				"Domain":     fixtures.PaymailDomain,
			})
	})

	t.Run("return paymails info for user (multiple paymails)", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForGivenUser(fixtures.UserWithMorePaymails)

		// when:
		res, _ := client.R().Get("/api/v1/paymails")

		// then:
		then.Response(res).
			IsOK().
			WithJSONTemplate(`{
			 "content": [
				{
				  "address": "{{.SecondPaymail.Address}}",
				  "alias": "{{.SecondPaymail.Alias}}",
				  "avatar": "/.*/",
				  "createdAt": "/.*/",
				  "deletedAt": null,
				  "domain": "{{.Domain}}",
				  "id": "/^[a-zA-Z0-9]{64}$/",
				  "metadata": "*",
				  "publicName": "{{.SecondPaymail.PublicName}}",
				  "updatedAt": "/.*/",
				  "xpubId": "{{.XPubID}}"
				},
				{
				  "address": "{{.FirstPaymail.Address}}",
				  "alias": "{{.FirstPaymail.Alias}}",
				  "avatar": "/.*/",
				  "createdAt": "/.*/",
				  "deletedAt": null,
				  "domain": "{{.Domain}}",
				  "id": "/^[a-zA-Z0-9]{64}$/",
				  "metadata": "*",
				  "publicName": "{{.FirstPaymail.PublicName}}",
				  "updatedAt": "/.*/",
				  "xpubId": "{{.XPubID}}"
				}
			 ],
			 "page": {
				"number": 1,
				"size": 50,
				"totalElements": 2,
				"totalPages": 1
			 }
			}`, map[string]any{
				"FirstPaymail": map[string]any{
					"Address":    strings.ToLower(fixtures.UserWithMorePaymails.Paymails[0]),
					"PublicName": fixtures.UserWithMorePaymails.Paymails[0],
					"Alias":      getAliasFromPaymail(t, fixtures.UserWithMorePaymails.Paymails[0]),
				},
				"SecondPaymail": map[string]any{
					"Address":    strings.ToLower(fixtures.UserWithMorePaymails.Paymails[1]),
					"PublicName": fixtures.UserWithMorePaymails.Paymails[1],
					"Alias":      getAliasFromPaymail(t, fixtures.UserWithMorePaymails.Paymails[1]),
				},
				"XPubID": fixtures.UserWithMorePaymails.XPubID(),
				"Domain": fixtures.PaymailDomain,
			})
	})

	t.Run("try to return paymails info for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v1/paymails")

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("return xpub info for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get("/api/v1/paymails")

		// then:
		then.Response(res).IsUnauthorized()
	})
}

func getAliasFromPaymail(t testing.TB, paymail string) (alias string) {
	parts := strings.SplitN(paymail, "@", 2)
	if len(parts) == 0 {
		t.Fatalf("Failed to parse paymail: %s", paymail)
	}
	alias = strings.ToLower(parts[0])
	return
}
