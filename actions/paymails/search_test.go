package paymails_test

import (
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
			WithJSONMatching(`{
			 "content": [
				{
				  "address": "{{.Address}}",
				  "alias": "{{.Alias}}",
				  "avatar": "{{ matchURL | orEmpty }}",
				  "createdAt": "{{ matchTimestamp }}",
				  "deletedAt": null,
				  "domain": "{{.Domain}}",
				  "id": "{{ matchID64 }}",
				  "metadata": "*",
				  "publicName": "{{.PublicName}}",
				  "updatedAt": "{{ matchTimestamp }}",
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
				"Address":    fixtures.Sender.DefaultPaymail(),
				"PublicName": fixtures.Sender.DefaultPaymail().PublicName(),
				"Alias":      fixtures.Sender.DefaultPaymail().Alias(),
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
			WithJSONMatching(`{
			 "content": [
				{
				  "address": "{{.SecondPaymail.Address}}",
				  "alias": "{{.SecondPaymail.Alias}}",
				  "avatar": "{{ matchURL | orEmpty }}",
				  "createdAt": "{{ matchTimestamp }}",
				  "deletedAt": null,
				  "domain": "{{.Domain}}",
				  "id": "{{ matchID64 }}",
				  "metadata": "*",
				  "publicName": "{{.SecondPaymail.PublicName}}",
				  "updatedAt": "{{ matchTimestamp }}",
				  "xpubId": "{{.XPubID}}"
				},
				{
				  "address": "{{.FirstPaymail.Address}}",
				  "alias": "{{.FirstPaymail.Alias}}",
				  "avatar": "{{ matchURL | orEmpty }}",
				  "createdAt": "{{ matchTimestamp }}",
				  "deletedAt": null,
				  "domain": "{{.Domain}}",
				  "id": "{{ matchID64 }}",
				  "metadata": "*",
				  "publicName": "{{.FirstPaymail.PublicName}}",
				  "updatedAt": "{{ matchTimestamp }}",
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
					"Address":    fixtures.UserWithMorePaymails.Paymails[0],
					"PublicName": fixtures.UserWithMorePaymails.Paymails[0].PublicName(),
					"Alias":      fixtures.UserWithMorePaymails.Paymails[0].Alias(),
				},
				"SecondPaymail": map[string]any{
					"Address":    fixtures.UserWithMorePaymails.Paymails[1],
					"PublicName": fixtures.UserWithMorePaymails.Paymails[1].PublicName(),
					"Alias":      fixtures.UserWithMorePaymails.Paymails[1].Alias(),
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
