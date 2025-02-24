package baseapi

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/api/manualtests"
	"github.com/bitcoin-sv/spv-wallet/api/manualtests/client"
	"github.com/stretchr/testify/require"
)

func TestBaseAPIRead(t *testing.T) {
	state := manualtests.NewState()
	err := state.Load()
	require.NoError(t, err)

	var calls = map[string]struct {
		client manualtests.ClientFactory
		call   func(*client.ClientWithResponses) (manualtests.Result, error)
	}{
		"sharedConfig[admin]": {
			client: manualtests.AdminClientFactory,
			call: func(c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.SharedConfigWithResponse(context.Background())
			},
		},
		"sharedConfig[user]": {
			client: manualtests.UserClientFactory,
			call: func(c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.SharedConfigWithResponse(context.Background())
			},
		},
	}
	for name, endpoint := range calls {
		t.Run(name, func(t *testing.T) {
			manualtests.APICallFor(t, endpoint.client).
				Call(endpoint.call).
				RequireSuccess()
		})
	}
}

func TestUnauthorized(t *testing.T) {
	state := manualtests.NewState()
	err := state.Load()
	require.NoError(t, err)

	var calls = map[string]struct {
		call func(*client.ClientWithResponses) (manualtests.Result, error)
	}{
		"sharedConfig": {
			call: func(c *client.ClientWithResponses) (manualtests.Result, error) {
				return c.SharedConfigWithResponse(context.Background())
			},
		},
	}

	for name, endpoint := range calls {
		t.Run(name+"_anonymous", func(t *testing.T) {
			manualtests.APICallForAnonymous(t).
				Call(endpoint.call).
				RequireUnauthorizedForAnonymous()
		})
	}

	for name, endpoint := range calls {
		t.Run(name+"_unknown_user", func(t *testing.T) {
			manualtests.APICallForUnknownUser(t).
				Call(endpoint.call).
				RequireUnauthorizedForUnknownUser()
		})
	}
}
