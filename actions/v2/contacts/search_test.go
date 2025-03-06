package contacts_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestSearchContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("No contacts", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().Get("/api/v2/contacts")

		// then:
		then.Response(res).
			HasStatus(200).
			WithJSONMatching(`
			{
				"content": [],
				"page": {
					"number": {{ .page }},
					"size": {{ .size }},
					"totalElements": {{ .totalElements }},
					"totalPages": {{ .totalPages }}
				}
			}`, map[string]any{
				"page":          1,
				"size":          0,
				"totalElements": 0,
				"totalPages":    0,
			})
	})

	t.Run("Search for all contacts", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		c1 := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)
		c3 := given.User(fixtures.Sender).HasContactTo(fixtures.UserWithMorePaymails)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().Get("/api/v2/contacts")

		// then:
		then.Response(res).
			HasStatus(200).
			WithJSONMatching(`
			{
				"content": [
				{
					"id": "{{ matchNumber }}",
					"createdAt": "{{ matchTimestamp }}",
					"updatedAt": "{{ matchTimestamp }}",
					"fullName": "{{ .fullName1 }}",
					"paymail": "{{ .paymail1 }}",
					"pubKey": "{{ matchHexWithLength 66 }}",
					"status": "{{ .status1 }}"
				},
				{
					"id": "{{ matchNumber }}",
					"createdAt": "{{ matchTimestamp }}",
					"updatedAt": "{{ matchTimestamp }}",
					"fullName": "{{ .fullName2 }}",
					"paymail": "{{ .paymail2 }}",
					"pubKey": "{{ matchHexWithLength 66 }}",
					"status": "{{ .status2 }}"
				}],
				"page": {
					"number": {{ .page }},
					"size": {{ .size }},
					"totalElements": {{ .totalElements }},
					"totalPages": {{ .totalPages }}
				}
			}`, map[string]any{
				"fullName1":     c3.FullName,
				"paymail1":      c3.Paymail,
				"status1":       c3.Status,
				"fullName2":     c1.FullName,
				"paymail2":      c1.Paymail,
				"status2":       c1.Status,
				"page":          1,
				"size":          2,
				"totalElements": 2,
				"totalPages":    1,
			})
	})

	t.Run("Search with user xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)
		given.User(fixtures.Sender).HasContactTo(fixtures.UserWithMorePaymails)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v2/contacts")

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrAdminAuthOnUserEndpoint.Code,
				"message": spverrors.ErrAdminAuthOnUserEndpoint.Message,
			})
	})

	t.Run("Search with pagination", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		c1 := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)
		given.User(fixtures.Sender).HasContactTo(fixtures.UserWithMorePaymails)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		pageQuery := fmt.Sprintf("page=%d&size=%d&sort=%s&sortBy=%s", 1, 1, "asc", "created_at")

		// when:
		res, _ := client.R().Get(fmt.Sprintf("/api/v2/contacts?%s", pageQuery))

		// then:
		then.Response(res).
			HasStatus(200).
			WithJSONMatching(`
			{
				"content": [
				{
					"id": "{{ matchNumber }}",
					"createdAt": "{{ matchTimestamp }}",
					"updatedAt": "{{ matchTimestamp }}",
					"fullName": "{{ .fullName }}",
					"paymail": "{{ .paymail }}",
					"pubKey": "{{ matchHexWithLength 66 }}",
					"status": "{{ .status }}"
				}
				],
				"page": {
					"number": {{ .page }},
					"size": {{ .size }},
					"totalElements": {{ .totalElements }},
					"totalPages": {{ .totalPages }}
				}
			}`, map[string]any{
				"fullName":      c1.FullName,
				"paymail":       c1.Paymail,
				"status":        c1.Status,
				"page":          1,
				"size":          1,
				"totalElements": 2,
				"totalPages":    2,
			})
	})

	t.Run("Search with conditions", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)
		c3 := given.User(fixtures.Sender).HasContactTo(fixtures.UserWithMorePaymails)
		given.User(fixtures.UserWithMorePaymails).HasContactTo(fixtures.Sender)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		conditionsQuery := fmt.Sprintf("paymail=%s", fixtures.UserWithMorePaymails.DefaultPaymail().String())

		// when:
		res, _ := client.R().Get(fmt.Sprintf("/api/v2/contacts?%s", conditionsQuery))

		// then:
		then.Response(res).
			HasStatus(200).
			WithJSONMatching(`
			{
				"content": [{
					"id": "{{ matchNumber }}",
					"createdAt": "{{ matchTimestamp }}",
					"updatedAt": "{{ matchTimestamp }}",
					"fullName": "{{ .fullName }}",
					"paymail": "{{ .paymail }}",
					"pubKey": "{{ matchHexWithLength 66 }}",
					"status": "{{ .status }}"
				}],
				"page": {
					"number": {{ .page }},
					"size": {{ .size }},
					"totalElements": {{ .totalElements }},
					"totalPages": {{ .totalPages }}
				}
			}`, map[string]any{
				"fullName":      c3.FullName,
				"paymail":       c3.Paymail,
				"status":        c3.Status,
				"page":          1,
				"size":          1,
				"totalElements": 1,
				"totalPages":    1,
			})
	})

	t.Run("Search with conditions and pagination", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)
		c3 := given.User(fixtures.Sender).HasContactTo(fixtures.UserWithMorePaymails)
		given.User(fixtures.UserWithMorePaymails).HasContactTo(fixtures.Sender)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		pageQuery := fmt.Sprintf("page=%d&size=%d&sort=%s&sortBy=%s", 1, 1, "asc", "created_at")
		conditionsQuery := fmt.Sprintf("&paymail=%s", fixtures.UserWithMorePaymails.DefaultPaymail().String())

		// when:
		res, _ := client.R().Get(fmt.Sprintf("/api/v2/contacts?%s%s", pageQuery, conditionsQuery))

		// then:
		then.Response(res).
			HasStatus(200).
			WithJSONMatching(`
			{
				"content": [
				{
					"id": "{{ matchNumber }}",
					"createdAt": "{{ matchTimestamp }}",
					"updatedAt": "{{ matchTimestamp }}",
					"fullName": "{{ .fullName }}",
					"paymail": "{{ .paymail }}",
					"pubKey": "{{ matchHexWithLength 66 }}",
					"status": "{{ .status }}"
				}],
				"page": {
					"number": {{ .page }},
					"size": {{ .size }},
					"totalElements": {{ .totalElements }},
					"totalPages": {{ .totalPages }}
				}
			}`, map[string]any{
				"fullName":      c3.FullName,
				"paymail":       c3.Paymail,
				"status":        c3.Status,
				"page":          1,
				"size":          1,
				"totalElements": 1,
				"totalPages":    1,
			})
	})
}
