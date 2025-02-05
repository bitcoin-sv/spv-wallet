package v2_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/config"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
)

func TestGETConfigsShared(t *testing.T) {
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(func(cfg *config.AppConfig) {
		cfg.Paymail.Domains = []string{"example.com"}
		cfg.ExperimentalFeatures.PikePaymentEnabled = true
		cfg.ExperimentalFeatures.PikeContactsEnabled = true
	}, testengine.WithV2())
	defer cleanup()

	t.Run("return shared config for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v2/configs/shared")

		// then:
		then.Response(res).
			IsOK().
			WithJSONf(`{
				"paymailDomains": ["example.com"],
				"experimentalFeatures": {
					"pikeContactsEnabled": true,
					"pikePaymentEnabled": true
				}
			}`)
	})

	t.Run("return shared config for user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().Get("/api/v2/configs/shared")

		// then:
		then.Response(res).
			IsOK().
			WithJSONf(`{
				"paymailDomains": ["example.com"],
				"experimentalFeatures": {
					"pikeContactsEnabled": true,
					"pikePaymentEnabled": true
				}
			}`)

	})

	t.Run("return unauthorized for anonymous requests", func(t *testing.T) {
		// given
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get("/api/v2/configs/shared")

		// then:
		then.Response(res).IsUnauthorized()
	})
}
