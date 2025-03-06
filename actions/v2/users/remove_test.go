package users_test

import (
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestCreateAndDeleteUser(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	userCandidate := fixtures.User{
		PrivKey: "xprv9s21ZrQH143K3QFk3G7fGtpKfi6ws96DVzeXpvvZLUafPBwnpfbX7A343GU9jMrbvkoJR3UrdCKjwPhXrPAmzDfQ8ipo3zLryFvj2ABH1hn",
		Paymails: []fixtures.Paymail{
			"test_user2@" + fixtures.PaymailDomain,
		},
	}
	publicKey := userCandidate.PublicKey().ToDERHex()

	var testState struct {
		userID string
	}

	t.Run("Create a user as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"publicKey": publicKey,
				"paymail": map[string]any{
					"address": userCandidate.DefaultPaymail(),
				},
			}).
			Post("/api/v2/admin/users")

		// then:
		then.Response(res).
			HasStatus(201).
			WithJSONMatching(`{
				"id": "{{ matchAddress }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"publicKey": "{{ .publicKey }}",
				"paymails": [
					{
						"alias": "{{ .alias }}",
						"avatar": "",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .paymail }}",
						"publicName": "{{ .publicName }}"
					}
				]
			}`, map[string]any{
				"publicKey":  publicKey,
				"paymail":    userCandidate.DefaultPaymail(),
				"publicName": userCandidate.DefaultPaymail().Alias(),
				"alias":      userCandidate.DefaultPaymail().Alias(),
			})

		// update:
		getter := then.Response(res).JSONValue()
		testState.userID = getter.GetString("id")
	})

	t.Run("Get new user by id as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetPathParam("id", testState.userID).
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
						"alias": "{{ .alias }}",
						"avatar": "",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .paymail }}",
						"publicName": "{{ .publicName }}"
					}
				]
			}`, map[string]any{
				"publicKey":  publicKey,
				"paymail":    userCandidate.DefaultPaymail(),
				"publicName": userCandidate.DefaultPaymail().Alias(),
				"alias":      userCandidate.DefaultPaymail().Alias(),
			})
	})

	t.Run("Delete user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForGivenUser(userCandidate)

		// when:
		res, _ := client.R().
			Delete("/api/v2/users/current")

		then.Response(res).IsOK()
	})

	t.Run("Try to get new user by id as admin after deletion", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetPathParam("id", testState.userID).
			Get("/api/v2/admin/users/{id}")

		then.Response(res).HasStatus(http.StatusNotFound).WithJSONf(apierror.ExpectedJSON("error-user-not-found", "user not found"))
	})

	t.Run("Try to make a request as deleted user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForGivenUser(userCandidate)

		// when:
		res, _ := client.R().Get("/api/v2/users/current")

		// then:
		then.Response(res).HasStatus(http.StatusUnauthorized).WithJSONf(apierror.ExpectedJSON("error-unauthorized", "unauthorized"))
	})
}
