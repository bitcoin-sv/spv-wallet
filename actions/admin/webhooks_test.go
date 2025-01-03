package admin_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
)

func TestAdminWebhooks(t *testing.T) {
	t.Run("subscribe, get and unsubscribe webhook", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithNotificationsEnabled())
		defer cleanup()

		// and:
		client := given.HttpClient().ForAdmin()

		// and:
		webhook := map[string]string{
			"url":         "http://localhost:8080",
			"tokenHeader": "Authorization",
			"tokenValue":  "123",
		}

		// when:
		res, _ := client.R().Get("/api/v1/admin/webhooks/subscriptions")

		// then:
		then.Response(res).
			IsOK().
			WithJSONf(`[]`)

		// when:
		res, _ = client.
			R().
			SetBody(webhook).
			Post("/api/v1/admin/webhooks/subscriptions")

		// then:
		then.Response(res).IsOK()

		// when:
		res, _ = client.R().Get("/api/v1/admin/webhooks/subscriptions")

		// then:
		then.Response(res).
			IsOK().
			WithJSONf(`[{
				"url": "http://localhost:8080",
				"banned": false
		}]`)

		// when:
		res, _ = client.
			R().
			SetBody(map[string]string{"url": webhook["url"]}).
			Delete("/api/v1/admin/webhooks/subscriptions")

		// then:
		then.Response(res).IsOK()

		// when:
		res, _ = client.R().Get("/api/v1/admin/webhooks/subscriptions")

		// then:
		then.Response(res).
			IsOK().
			WithJSONf(`[]`)
	})
}
